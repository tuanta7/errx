package errx

import (
	"net/http"
	"os"
	"sync"

	"github.com/tuanta7/errx/errors"
	"github.com/tuanta7/errx/parsers"
	"google.golang.org/grpc/codes"
)

const (
	DefaultMessage        = "internal server error"
	DefaultHTTPStatusCode = http.StatusInternalServerError
	DefaultGRPCStatusCode = uint32(codes.Internal)
)

type Errx struct {
	mu                 sync.RWMutex
	ErrorStatusCodeMap map[string]errors.StatusCode
	MessageMap         map[string]map[string]string // code -> language -> message
}

func NewWithDefaults() *Errx {
	// value types like sync.RWMutex don't need explicit initialization
	return &Errx{
		ErrorStatusCodeMap: errors.DefaultErrorStatusCodeMap,
		MessageMap:         make(map[string]map[string]string),
	}
}

func New() *Errx {
	return &Errx{
		ErrorStatusCodeMap: make(map[string]errors.StatusCode),
		MessageMap:         make(map[string]map[string]string),
	}
}

func (e *Errx) LoadMessages(filePath string, p parsers.Parser) error {
	bytesContents, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	m, err := p.Unmarshal(bytesContents)
	if err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	for code, messages := range m {
		e.MessageMap[code] = messages
	}

	return nil
}

func (e *Errx) RegisterMessage(code string, language, message string) {
	if _, exist := e.MessageMap[code]; !exist {
		e.MessageMap[code] = make(map[string]string)
	}

	e.MessageMap[code][language] = message
}

func (e *Errx) GetMessage(err *errors.Error, language string) string {
	if err == nil {
		return DefaultMessage
	}

	errorCode := err.Code()

	messages, exists := e.MessageMap[errorCode]
	if !exists {
		return getErrorMessage(err)
	}

	if localized, exists := messages[language]; exists {
		return localized
	}

	return getErrorMessage(err)
}

func getErrorMessage(err *errors.Error) string {
	if m := err.Message(); m != "" {
		return m
	}

	return DefaultMessage
}

func (e *Errx) RegisterErrorCode(code string, statusCode errors.StatusCode) {
	e.ErrorStatusCodeMap[code] = statusCode
}

func (e *Errx) RegisterHTTPErrorCode(code string, statusCode int) {
	if codeStatus, ok := e.ErrorStatusCodeMap[code]; ok {
		codeStatus.HTTPCode = statusCode
		return
	}

	e.ErrorStatusCodeMap[code] = errors.StatusCode{
		HTTPCode: statusCode,
	}
}

func (e *Errx) GetHTTPCode(err *errors.Error) int {
	if err == nil {
		return DefaultHTTPStatusCode
	}

	if code, ok := e.ErrorStatusCodeMap[err.Code()]; ok {
		return code.HTTPCode
	}

	return DefaultHTTPStatusCode
}

func (e *Errx) GetGRPCCode(err *errors.Error) uint32 {
	if err == nil {
		return DefaultGRPCStatusCode
	}

	if code, ok := e.ErrorStatusCodeMap[err.Code()]; ok {
		return code.GRPCCode
	}

	return DefaultGRPCStatusCode
}

func (e *Errx) GetHTTPResponse(err error, language string) (statusCode int, message string) {
	if err == nil {
		return http.StatusOK, ""
	}

	var errx *errors.Error
	if ok := errors.As(err, &errx); !ok {
		return e.GetHTTPCode(nil), e.GetMessage(nil, language)
	}

	statusCode = e.GetHTTPCode(errx)
	message = e.GetMessage(errx, language)
	return
}

func (e *Errx) GetGRPCResponse(err error, language string) (code uint32, message string) {
	if err == nil {
		return uint32(codes.OK), ""
	}

	var errx *errors.Error
	if ok := errors.As(err, &errx); !ok {
		return e.GetGRPCCode(nil), e.GetMessage(nil, language)
	}

	code = e.GetGRPCCode(errx)
	message = e.GetMessage(errx, language)
	return
}
