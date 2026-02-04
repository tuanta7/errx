package registry

import (
	"testing"

	"github.com/stretchr/testify/suite"
	errx2 "github.com/tuanta7/errx"
)

type MessageTestSuite struct {
	suite.Suite
	err        error
	wrappedErr *errx2.Error
	errx       *Registry
}

func (s *MessageTestSuite) SetupSuite() {
	s.err = errx2.New("test error")
	s.wrappedErr = errx2.New("wrapped message", s.err).WithCode("ERRX_TEST_ERROR")
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

func TestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}
