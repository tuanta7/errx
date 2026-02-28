package predefined

import (
	"github.com/tuanta7/errx"
)

var (
	ErrInternal                  = errx.New("internal error").WithCode("errx_internal")
	ErrServiceUnavailable        = errx.New("service unavailable").WithCode("errx_service_unavailable")
	ErrInvalidParameter          = errx.New("invalid parameter").WithCode("errx_invalid_parameter")
	ErrConnectionTimeout         = errx.New("connection timeout").WithCode("errx_connection_timeout")
	ErrOperationTimeout          = errx.New("operation timeout").WithCode("errx_operation_timeout")
	ErrRecordNotFound            = errx.New("record not found").WithCode("errx_record_not_found")
	ErrForeignKeyViolation       = errx.New("foreign key violation").WithCode("errx_foreign_key_violation")
	ErrUniqueConstraintViolation = errx.New("unique constraint violation").WithCode("errx_unique_constraint_violation")
)
