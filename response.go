package errx

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

func HTTPResponse(err *Error, language string) (statusCode int, message string) {
	if err == nil {
		return http.StatusOK, ""
	}

	statusCode = HTTPCode(err)
	message = Message(err.Code(), language)
	return
}

func GRPCResponse(err *Error, language string) (code uint32, message string) {
	if err == nil {
		return uint32(codes.OK), ""
	}

	code = GRPCCode(err)
	message = Message(err.Code(), language)
	return
}
