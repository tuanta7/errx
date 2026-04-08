package registry

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tuanta7/errx"
)

type MessageTestSuite struct {
	suite.Suite
	err        error
	wrappedErr *errx.Error
	errx       *Registry
}

func (s *MessageTestSuite) SetupSuite() {
	s.err = errx.New("test error")
	s.wrappedErr = errx.New("wrapped message", s.err).WithCode("ERRX_TEST_ERROR")
	s.errx = New()

	s.errx.RegisterMessage("ERRX_TEST_ERROR", "en", "test error")
	s.errx.RegisterMessage("ERRX_TEST_ERROR", "es", "error de prueba")
}

func (s *MessageTestSuite) TearDownSuite() {}

func (s *MessageTestSuite) TestErrorMessage() {
	s.Equal("test error", s.errx.GetMessage(s.wrappedErr, "en"))
	s.Equal("error de prueba", s.errx.GetMessage(s.wrappedErr, "es"))
	s.Equal("wrapped message", s.errx.GetMessage(s.wrappedErr, "vi"))
	s.Equal(DefaultMessage, s.errx.GetMessage(nil, "en"))
}

// TestRegionedTagFallback verifies that a regioned tag (e.g. "en-US") falls
// back to the registered base language ("en") when there is no exact match.
func (s *MessageTestSuite) TestRegionedTagFallback() {
	// "en-US" → falls back to "en"
	s.Equal("test error", s.errx.GetMessage(s.wrappedErr, "en-US"))
	// "es-MX" → falls back to "es"
	s.Equal("error de prueba", s.errx.GetMessage(s.wrappedErr, "es-MX"))
}

// TestMissingLanguageFallback verifies the full fallback chain when neither
// the exact tag nor the base language has a registered message.
func (s *MessageTestSuite) TestMissingLanguageFallback() {
	// "vi" is not registered → falls back to error's Message()
	s.Equal("wrapped message", s.errx.GetMessage(s.wrappedErr, "vi"))
	// "vi-VN" is not registered → falls back to error's Message()
	s.Equal("wrapped message", s.errx.GetMessage(s.wrappedErr, "vi-VN"))

	// Error with no registered code → falls back to error's Message()
	noCodeErr := errx.New("no code error")
	s.Equal("no code error", s.errx.GetMessage(noCodeErr, "en"))

	// nil error → DefaultMessage
	s.Equal(DefaultMessage, s.errx.GetMessage(nil, "en"))
}

func TestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}
