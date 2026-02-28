package registry

import (
	"errors"
	"net/http"
	"os"
	"sync"

	"github.com/tuanta7/errx"
	"github.com/tuanta7/errx/parsers"
	"google.golang.org/grpc/codes"
)

const (
	DefaultMessage        = "internal server error"
	DefaultHTTPStatusCode = http.StatusInternalServerError
	DefaultGRPCStatusCode = uint32(codes.Internal)
)

type Registry struct {
	mu            sync.RWMutex
	StatusCodeMap map[string]errx.StatusCode
	MessageMap    map[string]map[string]string // code -> language -> message
}

func New() *Registry {
	return &Registry{
		StatusCodeMap: make(map[string]errx.StatusCode),
		MessageMap:    make(map[string]map[string]string),
	}
}

func (e *Registry) LoadMessages(language, filePath string, p parsers.Parser) error {
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

	for code, message := range m {
		if _, exist := e.MessageMap[code]; !exist {
			e.MessageMap[code] = make(map[string]string)
		}
		e.MessageMap[code][language] = message
	}

	return nil
}

func (e *Registry) RegisterMessage(code, language, message string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exist := e.MessageMap[code]; !exist {
		e.MessageMap[code] = make(map[string]string)
	}

	e.MessageMap[code][language] = message
}

func (e *Registry) GetMessage(ex *errx.Error, language string) string {
	if ex == nil {
		return DefaultMessage
	}

	errorCode := ex.Code()

	e.mu.RLock()
	defer e.mu.RUnlock()

	messages, exists := e.MessageMap[errorCode]
	if !exists {
		return getErrorMessage(ex)
	}

	if localized, exists := messages[language]; exists {
		return localized
	}

	return getErrorMessage(ex)
}

func getErrorMessage(ex *errx.Error) string {
	if m := ex.Message(); m != "" {
		return m
	}

	return DefaultMessage
}

func (e *Registry) RegisterStatus(code string, statusCode errx.StatusCode) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.StatusCodeMap[code] = statusCode
}

func (e *Registry) RegisterHTTPStatus(code string, statusCode int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if codeStatus, ok := e.StatusCodeMap[code]; ok {
		codeStatus.HTTPCode = statusCode
		e.StatusCodeMap[code] = codeStatus
		return
	}

	e.StatusCodeMap[code] = errx.StatusCode{
		HTTPCode: statusCode,
	}
}

func (e *Registry) RegisterGRPCStatus(code string, statusCode uint32) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if codeStatus, ok := e.StatusCodeMap[code]; ok {
		codeStatus.GRPCCode = statusCode
		e.StatusCodeMap[code] = codeStatus
		return
	}

	e.StatusCodeMap[code] = errx.StatusCode{
		GRPCCode: statusCode,
	}
}

func (e *Registry) HTTPCode(err *errx.Error) int {
	if err == nil {
		return DefaultHTTPStatusCode
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	if code, ok := e.StatusCodeMap[err.Code()]; ok {
		return code.HTTPCode
	}

	return DefaultHTTPStatusCode
}

func (e *Registry) GRPCCode(ex *errx.Error) uint32 {
	if ex == nil {
		return DefaultGRPCStatusCode
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	if code, ok := e.StatusCodeMap[ex.Code()]; ok {
		return code.GRPCCode
	}

	return DefaultGRPCStatusCode
}

func (e *Registry) ResolveHTTP(err error, language string) (statusCode int, message string) {
	if err == nil {
		return http.StatusOK, ""
	}

	var ex *errx.Error
	if ok := errors.As(err, &ex); !ok {
		return e.HTTPCode(nil), e.GetMessage(nil, language)
	}

	statusCode = e.HTTPCode(ex)
	message = e.GetMessage(ex, language)
	return
}

func (e *Registry) ResolveGRPC(err error, language string) (code uint32, message string) {
	if err == nil {
		return uint32(codes.OK), ""
	}

	var ex *errx.Error
	if ok := errors.As(err, &ex); !ok {
		return e.GRPCCode(nil), e.GetMessage(nil, language)
	}

	code = e.GRPCCode(ex)
	message = e.GetMessage(ex, language)
	return
}
