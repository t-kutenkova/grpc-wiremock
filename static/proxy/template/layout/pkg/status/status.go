package statustocode

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

func GetCodeFromResponse(response *http.Response) codes.Code {
	if response == nil {
		return codes.Unknown
	}
	return codeFromHTTPStatus(response.StatusCode)
}

func GetStatusFromResponse(response *http.Response) int {
	if response == nil {
		return http.StatusNotFound
	}
	return response.StatusCode
}

func codeFromHTTPStatus(code int) codes.Code {
	switch code {
	case http.StatusContinue:
		return codes.OK
	case http.StatusSwitchingProtocols:
		return codes.OK
	case http.StatusProcessing:
		return codes.OK
	case http.StatusEarlyHints:
		return codes.OK
	case http.StatusOK:
		return codes.OK
	case http.StatusCreated:
		return codes.OK
	case http.StatusAccepted:
		return codes.OK
	case http.StatusNonAuthoritativeInfo:
		return codes.OK
	case http.StatusNoContent:
		return codes.OK
	case http.StatusResetContent:
		return codes.OK
	case http.StatusPartialContent:
		return codes.OK
	case http.StatusMultiStatus:
		return codes.OK
	case http.StatusAlreadyReported:
		return codes.OK
	case http.StatusIMUsed:
		return codes.OK
	case http.StatusMultipleChoices:
		return codes.OK
	case http.StatusMovedPermanently:
		return codes.OK
	case http.StatusFound:
		return codes.OK
	case http.StatusSeeOther:
		return codes.OK
	case http.StatusNotModified:
		return codes.OK
	case http.StatusUseProxy:
		return codes.OK
	case http.StatusTemporaryRedirect:
		return codes.OK
	case http.StatusPermanentRedirect:
		return codes.OK
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusPaymentRequired:
		return codes.PermissionDenied
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusMethodNotAllowed:
		return codes.InvalidArgument
	case http.StatusNotAcceptable:
		return codes.InvalidArgument
	case http.StatusProxyAuthRequired:
		return codes.Unauthenticated
	case http.StatusRequestTimeout:
		return codes.DeadlineExceeded
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusGone:
		return codes.NotFound
	case http.StatusLengthRequired:
		return codes.InvalidArgument
	case http.StatusPreconditionFailed:
		return codes.FailedPrecondition
	case http.StatusRequestEntityTooLarge:
		return codes.InvalidArgument
	case http.StatusRequestURITooLong:
		return codes.InvalidArgument
	case http.StatusUnsupportedMediaType:
		return codes.InvalidArgument
	case http.StatusRequestedRangeNotSatisfiable:
		return codes.OutOfRange
	case http.StatusExpectationFailed:
		return codes.InvalidArgument
	case http.StatusTeapot:
		return codes.OK
	case http.StatusMisdirectedRequest:
		return codes.InvalidArgument
	case http.StatusUnprocessableEntity:
		return codes.InvalidArgument
	case http.StatusLocked:
		return codes.PermissionDenied
	case http.StatusFailedDependency:
		return codes.FailedPrecondition
	case http.StatusTooEarly:
		return codes.Unavailable
	case http.StatusUpgradeRequired:
		return codes.InvalidArgument
	case http.StatusPreconditionRequired:
		return codes.FailedPrecondition
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusRequestHeaderFieldsTooLarge:
		return codes.InvalidArgument
	case http.StatusUnavailableForLegalReasons:
		return codes.PermissionDenied
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusBadGateway:
		return codes.Unavailable
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	case http.StatusHTTPVersionNotSupported:
		return codes.InvalidArgument
	case http.StatusVariantAlsoNegotiates:
		return codes.Internal
	case http.StatusInsufficientStorage:
		return codes.ResourceExhausted
	case http.StatusLoopDetected:
		return codes.Internal
	case http.StatusNotExtended:
		return codes.Unimplemented
	case http.StatusNetworkAuthenticationRequired:
		return codes.Unauthenticated
	}
	return codes.Unknown
}
