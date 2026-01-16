package errx

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MessageTestSuite struct {
	suite.Suite
	err  error
	errx *Error
}

func (s *MessageTestSuite) SetupSuite() {
	s.err = New("test error")
	s.errx = New("wrapped message", s.err).WithCode("ERRX_TEST_ERROR")
	RegisterMessage("ERRX_TEST_ERROR", "en", "test error")
	RegisterMessage("ERRX_TEST_ERROR", "es", "error de prueba")
}

func (s *MessageTestSuite) TearDownSuite() {}

func (s *MessageTestSuite) TestDefinedError() {
	s.Equal("test error", Message(s.errx, "en"))
	s.Equal("error de prueba", Message(s.errx, "es"))
	s.Equal("wrapped message", Message(s.errx, "vi"))
	s.Equal(DefaultMessage, Message(nil, "en"))
}

func TestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}
