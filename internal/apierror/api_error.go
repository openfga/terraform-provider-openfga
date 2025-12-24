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

type statusCoder interface {
	StatusCode() int
}

type responseStatusCoder interface {
	ResponseStatusCode() int
}

func IsStatusNotFound(err error) bool {
	if err == nil {
		return false
	}

	// Check for special case: 400 Bad Request with authorization_model_not_found error code
	var ve fgaValidationErr
	if errors.As(err, &ve) && ve.ResponseStatusCode() == http.StatusBadRequest {
		if ve.ResponseCode() == openfga.ERRORCODE_AUTHORIZATION_MODEL_NOT_FOUND {
			return true
		}
	}

	// Check for errors implementing ResponseStatusCode()
	var rsc responseStatusCoder
	if errors.As(err, &rsc) && rsc.ResponseStatusCode() == http.StatusNotFound {
		return true
	}

	// Check for errors implementing StatusCode()
	var sc statusCoder
	if errors.As(err, &sc) && sc.StatusCode() == http.StatusNotFound {
		return true
	}

	// Fallback to existing check for FgaApiNotFoundError
	var nf openfga.FgaApiNotFoundError
	if errors.As(err, &nf) && nf.ResponseStatusCode() == http.StatusNotFound {
		return true
	}

	return false
}

func IsExpectedOneResultError(err error) bool {
	return errors.Is(err, ErrNotExactlyOne)
}
