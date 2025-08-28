package apierror

import (
	"errors"
	"net/http"

	openfga "github.com/openfga/go-sdk"
)

var ErrNotExactlyOne = errors.New("expected one result")

type fgaValidationErr interface {
	ResponseStatusCode() int
	ResponseCode() openfga.ErrorCode
}

func IsStatusNotFound(err error) bool {
	if err == nil {
		return false
	}

	var ve fgaValidationErr
	if errors.As(err, &ve) && ve.ResponseStatusCode() == http.StatusBadRequest {
		if ve.ResponseCode() == openfga.ERRORCODE_AUTHORIZATION_MODEL_NOT_FOUND {
			return true
		}
	}

	var nf openfga.FgaApiNotFoundError
	if errors.As(err, &nf) && nf.ResponseStatusCode() == http.StatusNotFound {
		return true
	}

	return false
}

func IsExpectedOneResultError(err error) bool {
	return errors.Is(err, ErrNotExactlyOne)
}
