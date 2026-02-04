package registry

import (
	"errors"
	"net/http"
	"os"
	"sync"

	x "github.com/tuanta7/errx"
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
	StatusCodeMap map[string]x.StatusCode
	MessageMap    map[string]map[string]string // code -> language -> message
}

func New() *Registry {
	return &Registry{
		StatusCodeMap: make(map[string]x.StatusCode),
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

func (e *Registry) RegisterMessage(code string, language, message string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exist := e.MessageMap[code]; !exist {
		e.MessageMap[code] = make(map[string]string)
	}

	e.MessageMap[code][language] = message
}

func (e *Registry) GetMessage(errx *x.Error, language string) string {
	if errx == nil {
		return DefaultMessage
	}

	errorCode := errx.Code()

	e.mu.RLock()
	defer e.mu.RUnlock()

	messages, exists := e.MessageMap[errorCode]
	if !exists {
		return getErrorMessage(errx)
	}

	if localized, exists := messages[language]; exists {
		return localized
	}

	return getErrorMessage(errx)
}

func getErrorMessage(errx *x.Error) string {
	if m := errx.Message(); m != "" {
		return m
	}

	return DefaultMessage
}

func (e *Registry) RegisterStatus(code string, statusCode x.StatusCode) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.StatusCodeMap[code] = statusCode
}

func (e *Registry) RegisterHTTPStatus(code string, statusCode int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if codeStatus, ok := e.StatusCodeMap[code]; ok {
		codeStatus.HTTPCode = statusCode
		return
	}

	e.StatusCodeMap[code] = x.StatusCode{
		HTTPCode: statusCode,
	}
}

func (e *Registry) RegisterGRPCStatus(code string, statusCode uint32) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if codeStatus, ok := e.StatusCodeMap[code]; ok {
		codeStatus.GRPCCode = statusCode
		return
	}

	e.StatusCodeMap[code] = x.StatusCode{
		GRPCCode: statusCode,
	}
}

func (e *Registry) HTTPCode(err *x.Error) int {
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

func (e *Registry) GRPCCode(errx *x.Error) uint32 {
	if errx == nil {
		return DefaultGRPCStatusCode
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	if code, ok := e.StatusCodeMap[errx.Code()]; ok {
		return code.GRPCCode
	}

	return DefaultGRPCStatusCode
}

func (e *Registry) ResolveHTTP(err error, language string) (statusCode int, message string) {
	if err == nil {
		return http.StatusOK, ""
	}

	var errx *x.Error
	if ok := errors.As(err, &errx); !ok {
		return e.HTTPCode(nil), e.GetMessage(nil, language)
	}

	statusCode = e.HTTPCode(errx)
	message = e.GetMessage(errx, language)
	return
}

func (e *Registry) ResolveGRPC(err error, language string) (code uint32, message string) {
	if err == nil {
		return uint32(codes.OK), ""
	}

	var errx *x.Error
	if ok := errors.As(err, &errx); !ok {
		return e.GRPCCode(nil), e.GetMessage(nil, language)
	}

	code = e.GRPCCode(errx)
	message = e.GetMessage(errx, language)
	return
}
