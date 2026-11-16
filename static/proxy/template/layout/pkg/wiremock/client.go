package wiremock

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	statustocode "grpc-proxy/pkg/status"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var client = http.Client{Timeout: time.Second * 3}

const (
	streamCursor = "streamCursor"
	streamSize   = "streamSize"
)

// skippedResponseHeaders are HTTP response headers that must not be forwarded
// to gRPC response metadata: hop-by-hop headers, headers describing the HTTP
// transport/body and WireMock service headers.
var skippedResponseHeaders = map[string]struct{}{
	// hop-by-hop
	"connection":          {},
	"keep-alive":          {},
	"proxy-authenticate":  {},
	"proxy-authorization": {},
	"proxy-connection":    {},
	"te":                  {},
	"trailer":             {},
	"transfer-encoding":   {},
	"upgrade":             {},
	// transport/body
	"content-type":     {},
	"content-length":   {},
	"content-encoding": {},
	"date":             {},
	"server":           {},
	"vary":             {},
	// wiremock
	"matched-stub-id":   {},
	"matched-stub-name": {},
	"streamsize":        {},
}

func DefaultRequest(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return httpRequest, nil
}

func RequestWithCursor(ctx context.Context, url string, cursor int, body io.Reader) (*http.Request, error) {
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Add(streamCursor, strconv.Itoa(cursor))
	return httpRequest, nil
}

func DoRequestDefault(request *http.Request) ([]byte, error) {
	response, err := doRequest(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, status.Error(http.StatusBadGateway, fmt.Sprintf("read http response: %v", err))
	}

	return body, nil
}

// DoRequestUnary returns the response body and the response headers that
// should be forwarded to the gRPC client as header metadata.
func DoRequestUnary(request *http.Request) ([]byte, metadata.MD, error) {
	response, err := doRequest(request)
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, status.Error(http.StatusBadGateway, fmt.Sprintf("read http response: %v", err))
	}

	return body, headersToMetaData(response.Header), nil
}

func DoRequestWithStreamSize(request *http.Request) ([]byte, int, error) {
	response, err := doRequest(request)
	if err != nil {
		return nil, 0, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, 0, status.Error(http.StatusBadGateway, fmt.Sprintf("read http response: %v", err))
	}

	streamSizeRaw := response.Header.Get(streamSize)
	if len(streamSizeRaw) == 0 {
		return nil, 0, status.Error(http.StatusBadGateway, fmt.Sprintf("read http response, stream size: %v", err))
	}

	streamSize, err := strconv.Atoi(streamSizeRaw)
	if err != nil {
		return nil, 0, fmt.Errorf("convert streamSize header to int: %w", err)
	}

	return body, streamSize, nil
}

func enrichWithMetaData(request *http.Request) *http.Request {
	if md, ok := metadata.FromIncomingContext(request.Context()); ok {
		for mdKey, mdValue := range md {
			if len(mdValue) > 0 {
				if mdKey == ":authority" {
					request.Host = mdValue[0]
				} else {
					request.Header.Set(mdKey, mdValue[0])
				}
			}
		}
	}

	return request
}

func headersToMetaData(header http.Header) metadata.MD {
	skipped := make(map[string]struct{})
	for _, value := range header.Values("Connection") {
		for _, name := range strings.Split(value, ",") {
			skipped[strings.ToLower(strings.TrimSpace(name))] = struct{}{}
		}
	}

	md := metadata.MD{}
	for name, values := range header {
		key := strings.ToLower(name)
		if _, ok := skippedResponseHeaders[key]; ok {
			continue
		}
		if _, ok := skipped[key]; ok {
			continue
		}
		if strings.HasPrefix(key, "grpc-") || !isValidMetaDataKey(key) {
			continue
		}
		md.Append(key, values...)
	}

	return md
}

// isValidMetaDataKey reports whether key is allowed as gRPC metadata key: [0-9a-z-_.].
func isValidMetaDataKey(key string) bool {
	if key == "" {
		return false
	}
	for _, r := range key {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '.') {
			return false
		}
	}
	return true
}

func doRequest(request *http.Request) (*http.Response, error) {
	request = enrichWithMetaData(request)

	httpResponse, err := client.Do(request)
	if err != nil {
		code := statustocode.GetCodeFromResponse(httpResponse)
		return nil, status.Error(code, fmt.Sprintf("do http response: %v", err))
	}

	if httpStatus := statustocode.GetStatusFromResponse(httpResponse); httpStatus >= http.StatusBadRequest {
		code := statustocode.GetCodeFromResponse(httpResponse)
		message := fmt.Sprintf("wiremock bad status: %d", httpStatus)
		if body, err := io.ReadAll(httpResponse.Body); err == nil {
			message += "\n" + string(body)
		}
		return nil, status.Error(code, message)
	}

	return httpResponse, nil
}
