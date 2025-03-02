package errs_test

import (
	"bytes"
	"net/http"
	"testing"

	errs "github.com/Lucasvmarangoni/logella/err"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const (
	message = "teste message"
	cause   = "test error"
)

func handler(code int) (string, error) {
	_, err := service(code)
	return "", errs.Trace(err)
}

func service(code int) (string, error) {
	err := repository(code)
	if err != nil {
		return "", errs.Trace(err)
	}
	return "", nil
}

func repository(code int) error {
	err := externalMethod()
	return errs.Wrap(err, code).Msg(message)
}

func externalMethod() error {
	return errors.New(cause)
}

func TestErrorOutput(t *testing.T) {

	t.Run("(ToClient) Should have the same error message and error code", func(t *testing.T) {
		expectedCode := http.StatusInternalServerError
		_, err := handler(expectedCode)

		if err == nil {
			t.Fatal("Expected an error, but got nil. The handler did not return any error as expected.")
		}

		expectedCause := http.StatusText(http.StatusInternalServerError)
		errCause := errs.Unwrap(err).ToClient().Error()
		assert.Equal(t, expectedCause, errCause,
			"Expected error message '%s', but got '%s'. This indicates a mismatch in the error handling or messaging.",
			expectedCause, errCause)

		errCode := errs.Unwrap(err).Code
		assert.Equal(t, expectedCode, errCode,
			"Expected error code %d, but got %d. This suggests an issue with the status code returned by the handler.",
			expectedCode, errCode)

		expectedMessage := message
		errMessage := errs.Unwrap(err).Message
		assert.Equal(t, expectedMessage, errMessage,
			"Expected error message %d, but got %d. This suggests an issue with the error message returned by the handler.",
			expectedMessage, errMessage)
	})

	t.Run("(Stack) Should have the same error message and error code", func(t *testing.T) {
		expectedCode := http.StatusBadRequest
		_, err := handler(expectedCode)

		if err == nil {
			t.Fatal("Expected an error, but got nil. The handler did not return any error as expected.")
		}

		expectedCause := cause
		errCause := errs.Unwrap(err).ToClient().Error()
		assert.Equal(t, expectedCause, errCause,
			"Expected error message '%s', but got '%s'. This indicates a mismatch in the error handling or messaging.",
			expectedCause, errCause)

		errCode := errs.Unwrap(err).Code
		assert.Equal(t, expectedCode, errCode,
			"Expected error code %d, but got %d. This suggests an issue with the status code returned by the handler.",
			expectedCode, errCode)

		expectedMessage := message
		errMessage := errs.Unwrap(err).Message
		assert.Equal(t, expectedMessage, errMessage,
			"Expected error message %d, but got %d. This suggests an issue with the error message returned by the handler.",
			expectedMessage, errMessage)
	})
}

func TestLogOutput(t *testing.T) {
	t.Run("(Stack) Log should have the same error message, error code, path and trace", func(t *testing.T) {
		expectedCode := http.StatusInternalServerError
		_, err := handler(expectedCode)

		var buf bytes.Buffer
		log.Logger = zerolog.New(&buf)

		log.Error().Err(errs.Unwrap(err).Stack()).Msg("example error")
		logOutput := buf.String()

		assert.Contains(t, logOutput, `example error`)
		assert.Contains(t, logOutput, cause)
		assert.Contains(t, logOutput, `path errs_test.go:35 trace repository ➤ service ➤ handler","message":"example error"`)
	})

	t.Run("(ToClient) Log should have the same error message and 'Internal Server Error' cause and not have the real cause, path and trace", func(t *testing.T) {
		expectedCode := http.StatusInternalServerError
		_, err := handler(expectedCode)

		var buf bytes.Buffer
		log.Logger = zerolog.New(&buf)

		log.Error().Err(errs.Unwrap(err).ToClient()).Msg("example error")
		logOutput := buf.String()

		assert.Contains(t, logOutput, `example error`)
		assert.Contains(t, logOutput, http.StatusText(expectedCode))

		assert.NotContains(t, logOutput, cause)
		assert.NotContains(t, logOutput, `path errs_test.go:35 trace repository ➤ service ➤ handler","message":"example error"`)
	})
}
