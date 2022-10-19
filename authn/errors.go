package authn

import (
	"net/http"

	"github.com/aserto-dev/errors"
	"google.golang.org/grpc/codes"
)

var (
	// Unknown error ID. It's returned when the implementation has not returned another AsertoError.
	ErrUnknown = newErr("E40000", codes.Internal, http.StatusInternalServerError, "an unknown error has occurred")
	// Returned when authentication has failed or is not possible
	ErrAuthenticationFailed = newErr("E40001", codes.FailedPrecondition, http.StatusUnauthorized, "authentication failed")
)

func newErr(code string, statusCode codes.Code, httpCode int, msg string) *errors.AsertoError {
	return errors.NewAsertoError(code, statusCode, httpCode, msg)
}
