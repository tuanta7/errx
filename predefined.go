package errx

var (
	ErrInternal          = New("internal error").WithCode("errx_internal")
	ErrInvalidParameter  = New("invalid parameter").WithCode("errx_invalid_parameter")
	ErrConnectionTimeout = New("connection timeout").WithCode("errx_connection_timeout")
	ErrOperationTimeout  = New("operation timeout").WithCode("errx_operation_timeout")
	ErrRecordNotFound    = New("record not found").WithCode("errx_record_not_found")
)

var (
	ErrForeignKeyViolation       = New("foreign key violation").WithCode("errx_foreign_key_violation")
	ErrUniqueConstraintViolation = New("unique constraint violation").WithCode("errx_unique_constraint_violation")
)
