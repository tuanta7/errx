package predefined

import (
	"net/http"

	"github.com/tuanta7/errx"
	"google.golang.org/grpc/codes"
)

var (
	DefaultErrorStatusCodeMap = map[string]errx.StatusCode{
		ErrInternal.Code():                  {http.StatusInternalServerError, uint32(codes.Internal)},
		ErrServiceUnavailable.Code():        {http.StatusServiceUnavailable, uint32(codes.Unavailable)},
		ErrInvalidParameter.Code():          {http.StatusBadRequest, uint32(codes.InvalidArgument)},
		ErrConnectionTimeout.Code():         {http.StatusServiceUnavailable, uint32(codes.Unavailable)},
		ErrOperationTimeout.Code():          {http.StatusGatewayTimeout, uint32(codes.DeadlineExceeded)},
		ErrRecordNotFound.Code():            {http.StatusNotFound, uint32(codes.NotFound)},
		ErrForeignKeyViolation.Code():       {http.StatusBadRequest, uint32(codes.InvalidArgument)},
		ErrUniqueConstraintViolation.Code(): {http.StatusBadRequest, uint32(codes.AlreadyExists)},
	}
)
