package errx

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

type StatusCode struct {
	HTTPCode int
	GRPCCode uint32
}

var globalErrorStatusCodeMap = map[string]StatusCode{
	ErrInternal.Code():                  {http.StatusInternalServerError, uint32(codes.Internal)},
	ErrInvalidParameter.Code():          {http.StatusBadRequest, uint32(codes.InvalidArgument)},
	ErrConnectionTimeout.Code():         {http.StatusServiceUnavailable, uint32(codes.Unavailable)},
	ErrOperationTimeout.Code():          {http.StatusGatewayTimeout, uint32(codes.DeadlineExceeded)},
	ErrRecordNotFound.Code():            {http.StatusNotFound, uint32(codes.NotFound)},
	ErrForeignKeyViolation.Code():       {http.StatusBadRequest, uint32(codes.InvalidArgument)},
	ErrUniqueConstraintViolation.Code(): {http.StatusBadRequest, uint32(codes.AlreadyExists)},
}

func RegisterErrorCode(code string, statusCode StatusCode) {
	globalErrorStatusCodeMap[code] = statusCode
}

func RegisterHTTPErrorCode(code string, statusCode int) {
	if codeStatus, ok := globalErrorStatusCodeMap[code]; ok {
		codeStatus.HTTPCode = statusCode
		return
	}

	globalErrorStatusCodeMap[code] = StatusCode{statusCode, 0}
}

func HTTPCode(err *Error) int {
	if err == nil {
		return http.StatusInternalServerError
	}

	if code, ok := globalErrorStatusCodeMap[err.Code()]; ok {
		return code.HTTPCode
	}

	return http.StatusInternalServerError
}

func GRPCCode(err *Error) uint32 {
	if err == nil {
		return uint32(codes.Internal)
	}

	if code, ok := globalErrorStatusCodeMap[err.Code()]; ok {
		return code.GRPCCode
	}

	return uint32(codes.Internal)
}
