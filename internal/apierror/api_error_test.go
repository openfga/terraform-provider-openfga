package apierror

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	openfga "github.com/openfga/go-sdk"
)

func TestIsExpectedOneResultError(t *testing.T) {
	t.Run("wrapped with %%w", func(t *testing.T) {
		err := fmt.Errorf("%w but received: %d", ErrNotExactlyOne, 0)
		if !IsExpectedOneResultError(err) {
			t.Fatalf("expected true for wrapped error")
		}
	})

	t.Run("nil error", func(t *testing.T) {
		var err error
		if IsExpectedOneResultError(err) {
			t.Fatalf("expected false for nil error")
		}
	})

	t.Run("unrelated error", func(t *testing.T) {
		err := errors.New("some other error")
		if IsExpectedOneResultError(err) {
			t.Fatalf("expected false for unrelated error")
		}
	})
}

func TestHandleAPIError(t *testing.T) {
	testCases := []struct {
		name          string
		givenErr      error
		expectedValue bool
	}{
		{
			name:          "returns true for 404 not found",
			givenErr:      createNotFoundError(),
			expectedValue: true,
		},
		{
			name:          "returns false for generic 400 api error",
			givenErr:      createBadRequestError(),
			expectedValue: false,
		},
		{
			name:          "returns false for non-FGA error",
			givenErr:      fmt.Errorf("400"),
			expectedValue: false,
		},
		{
			name:          "returns true for 400 validation error: authorization_model_not_found",
			givenErr:      createValidationAuthModelNotFoundError(),
			expectedValue: true,
		},
		{
			name:          "returns false for 400 validation error with different code",
			givenErr:      createValidationOtherCodeError(),
			expectedValue: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsStatusNotFound(tc.givenErr)
			if result != tc.expectedValue {
				t.Errorf("expected %v, but got %v", tc.expectedValue, result)
			}
		})
	}
}

func createHttpErrorResponse(statusCode int) *http.Response {
	req, _ := http.NewRequest("GET", "https://api.fga.example/stores/test-store/check", nil)
	resp := &http.Response{
		StatusCode: statusCode,
		Header:     make(http.Header),
		Request:    req,
	}
	resp.Header.Set("Fga-Request-Id", "test-request-id")
	return resp
}

func createNotFoundError() error {
	resp := createHttpErrorResponse(http.StatusNotFound)
	return openfga.NewFgaApiNotFoundError(
		"check",
		map[string]any{"test": "body"},
		resp,
		[]byte(`{"error":"not found"}`),
		"test-store-id",
	)
}

func createBadRequestError() error {
	resp := createHttpErrorResponse(http.StatusBadRequest)
	return openfga.NewFgaApiError(
		"check",
		map[string]any{"test": "body"},
		resp,
		[]byte(`{"error":"bad request"}`),
		"test-store-id",
	)
}

type validationErr struct {
	resp *http.Response
	code openfga.ErrorCode
	msg  string
}

func (e validationErr) Error() string                   { return e.msg }
func (e validationErr) ResponseStatusCode() int         { return e.resp.StatusCode }
func (e validationErr) ResponseCode() openfga.ErrorCode { return e.code }

func createValidationAuthModelNotFoundError() error {
	return validationErr{
		resp: createHttpErrorResponse(http.StatusBadRequest),
		code: openfga.ERRORCODE_AUTHORIZATION_MODEL_NOT_FOUND,
		msg:  "authorization_model_not_found",
	}
}

func createValidationOtherCodeError() error {
	return validationErr{
		resp: createHttpErrorResponse(http.StatusBadRequest),
		code: openfga.ERRORCODE_TYPE_NOT_FOUND,
		msg:  "type_not_found",
	}
}
