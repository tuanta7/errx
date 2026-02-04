package errx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ErrorTypeTestSuite struct {
	suite.Suite
	err  error
	errx *Error
}

func (s *ErrorTypeTestSuite) SetupSuite() {
	s.err = errors.New("test error")
	s.errx = New("wrapped message", s.err)
}

func (s *ErrorTypeTestSuite) TearDownSuite() {}

func (s *ErrorTypeTestSuite) TestGetters() {
	errxWithCode := s.errx.WithCode("errx_test_code")
	s.Equal("errx_test_code", errxWithCode.Code())
	s.Equal("wrapped message", s.errx.Message())
}

func (s *ErrorTypeTestSuite) TestErrorUnwrap() {
	s.Equal(s.err, s.errx.Unwrap())
}

func (s *ErrorTypeTestSuite) TestErrorIsMethod() {
	// Is() is used for error comparison
	s.False(errors.Is(s.err, s.errx))
	s.True(errors.Is(s.errx, s.err))
	s.False(errors.Is(s.errx, New("test error")))
	s.False(errors.Is(s.errx, New("another error")))
}

func (s *ErrorTypeTestSuite) TestErrorAsMethod() {
	var e error = s.errx

	// As() is used for type assertion
	s.False(errors.As(s.err, &s.errx))
	s.True(errors.As(e, &s.errx))
}

func (s *ErrorTypeTestSuite) TestErrorFormat() {
	s.Equal("test error", s.errx.Error())
	s.Equal("wrapped message", s.errx.Message())
}

func TestClientRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorTypeTestSuite))
}
