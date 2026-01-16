package errx

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
)

func HTTPResponse(err error, language string) (statusCode int, message string) {
	if err == nil {
		return http.StatusOK, ""
	}

	var errx *Error
	ok := errors.As(err, &errx)
	if !ok {
		return HTTPCode(nil), Message("", language)
	}

	statusCode = HTTPCode(errx)
	message = Message(errx.Code(), language)
	return
}

func GRPCResponse(err error, language string) (code uint32, message string) {
	if err == nil {
		return uint32(codes.OK), ""
	}

	var errx *Error
	ok := errors.As(err, &errx)
	if !ok {
		return GRPCCode(nil), Message("", language)
	}

	code = GRPCCode(errx)
	message = Message(errx.Code(), language)
	return
}
