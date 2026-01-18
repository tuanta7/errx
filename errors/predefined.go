package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

var (
	ErrInternal                  = New("internal error").WithCode("errx_internal")
	ErrServiceUnavailable        = New("service unavailable").WithCode("errx_service_unavailable")
	ErrInvalidParameter          = New("invalid parameter").WithCode("errx_invalid_parameter")
	ErrConnectionTimeout         = New("connection timeout").WithCode("errx_connection_timeout")
	ErrOperationTimeout          = New("operation timeout").WithCode("errx_operation_timeout")
	ErrRecordNotFound            = New("record not found").WithCode("errx_record_not_found")
	ErrForeignKeyViolation       = New("foreign key violation").WithCode("errx_foreign_key_violation")
	ErrUniqueConstraintViolation = New("unique constraint violation").WithCode("errx_unique_constraint_violation")
)

var (
	DefaultErrorStatusCodeMap = map[string]StatusCode{
		ErrInternal.Code():                  {http.StatusInternalServerError, uint32(codes.Internal)},
		ErrServiceUnavailable.Code():        {http.StatusServiceUnavailable, uint32(codes.Unavailable)},
		ErrInvalidParameter.Code():          {http.StatusBadRequest, uint32(codes.InvalidArgument)},
		ErrConnectionTimeout.Code():         {http.StatusServiceUnavailable, uint32(codes.Unavailable)},
		ErrOperationTimeout.Code():          {http.StatusGatewayTimeout, uint32(codes.DeadlineExceeded)},
		ErrRecordNotFound.Code():            {http.StatusNotFound, uint32(codes.NotFound)},
		ErrForeignKeyViolation.Code():       {http.StatusBadRequest, uint32(codes.InvalidArgument)},
		ErrUniqueConstraintViolation.Code(): {http.StatusBadRequest, uint32(codes.AlreadyExists)},
	}
)
