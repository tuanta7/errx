package errors

type Error struct {
	code    string
	message string
	err     error
}

func New(message string, err ...error) *Error {
	e := &Error{
		message: message,
	}

	if len(err) > 0 {
		e.err = err[0]
	}

	return e
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	return e.message
}

func (e *Error) WithCode(code string) *Error {
	c := *e
	c.code = code
	return &c
}

func (e *Error) Code() string {
	return e.code
}

func (e *Error) Message() string {
	return e.message
}
