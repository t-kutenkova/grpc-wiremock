package wiremock

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/stretchr/testify/require"
)

func Test_enrichWithMetaData(t *testing.T) {
	tests := []struct {
		name string
		md   metadata.MD

		wantHost   string
		wantHeader http.Header
	}{
		{
			name:       "empty metadata means empty host and header",
			md:         metadata.MD{},
			wantHost:   "",
			wantHeader: http.Header{},
		},
		{
			name: "authority metadata fills request host",
			md: metadata.MD{
				":authority": []string{
					"push-sender.test",
				},
			},
			wantHost:   "push-sender.test",
			wantHeader: http.Header{},
		},
		{
			name: "custom metadata fills request header",
			md: metadata.MD{
				"custom": []string{"testHeader"},
				"content-type": []string{
					"application/grpc",
				},
			},
			wantHost: "",
			wantHeader: http.Header{
				"Custom":       []string{"testHeader"},
				"Content-Type": []string{"application/grpc"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := DefaultRequest(
				metadata.NewIncomingContext(context.Background(), tt.md),
				"/Notify", bytes.NewReader([]byte{}))
			require.NoError(t, err)

			request = enrichWithMetaData(request)

			require.Equal(t, tt.wantHost, request.Host)
			require.Equal(t, tt.wantHeader, request.Header)
		})
	}
}

func Test_headersToMetaData(t *testing.T) {
	tests := []struct {
		name   string
		header http.Header
		wantMD metadata.MD
	}{
		{
			name:   "empty header means empty metadata",
			header: http.Header{},
			wantMD: metadata.MD{},
		},
		{
			name: "custom headers are forwarded in lower case",
			header: http.Header{
				"Content-Disposition": []string{`attachment; filename="invoice.pdf"`},
				"X-Custom":            []string{"a", "b"},
			},
			wantMD: metadata.MD{
				"content-disposition": []string{`attachment; filename="invoice.pdf"`},
				"x-custom":            []string{"a", "b"},
			},
		},
		{
			name: "hop-by-hop, service and grpc headers are skipped",
			header: http.Header{
				"Connection":        []string{"keep-alive, X-Hop"},
				"Keep-Alive":        []string{"timeout=5"},
				"Transfer-Encoding": []string{"chunked"},
				"X-Hop":             []string{"hop"},
				"Content-Type":      []string{"application/json"},
				"Content-Length":    []string{"10"},
				"Date":              []string{"Mon, 16 Nov 2026 00:00:00 GMT"},
				"Matched-Stub-Id":   []string{"id"},
				"Streamsize":        []string{"3"},
				"Grpc-Status":       []string{"0"},
				"X-Bad:Key":         []string{"bad"},
				"X-Kept":            []string{"kept"},
			},
			wantMD: metadata.MD{
				"x-kept": []string{"kept"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantMD, headersToMetaData(tt.header))
		})
	}
}

func Test_DoRequestUnary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", `attachment; filename="invoice.pdf"`)
		_, _ = w.Write([]byte(`{"data":"cGRm"}`))
	}))
	defer server.Close()

	request, err := DefaultRequest(context.Background(), server.URL, bytes.NewReader([]byte{}))
	require.NoError(t, err)

	body, md, err := DoRequestUnary(request)
	require.NoError(t, err)
	require.Equal(t, `{"data":"cGRm"}`, string(body))
	require.Equal(t, metadata.MD{
		"content-disposition": []string{`attachment; filename="invoice.pdf"`},
	}, md)
}
