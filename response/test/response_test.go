package response_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	errs "github.com/Lucasvmarangoni/logella/err"
	"github.com/Lucasvmarangoni/logella/response"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestSendResponse_ReturnsErrorResponseOnNonErrsError(t *testing.T) {
	resp := response.New(errors.New("plain error"))

	assert.Equal(t, "Internal Server Error", resp.Status)
	assert.Equal(t, "Invalid error type passed to response.New: expected *errs.Error. Use errs.Wrap.", resp.Error)
	assert.Nil(t, resp.Err) 
}


func TestSendResponse_CorrectJSONAndLog(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = log.Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = &buf
		w.TimeFormat = time.RFC3339
	}))
	recorder := httptest.NewRecorder()

	fixed := time.Date(2025, 5, 15, 17, 10, 5, 0, time.FixedZone("-03:00", -3*60*60))

	errTest := errors.New("some error")
	e := errs.Wrap(errTest, http.StatusBadRequest)
	response.New(e).
		Req("123").
		User("12345").
		Log("LOG MESSAGE").
		Date(&fixed).
		Send(recorder)

	result := recorder.Result()
	defer result.Body.Close()

	res := `{"error":"some error","message":"some error | ","status":"Bad Request","request_id":"123","user_id":"12345","timestamp":"2025-05-15T17:10:05-03:00"}`

	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	bodyBytes, err := io.ReadAll(result.Body)
	assert.NoError(t, err)

	bodyStr := string(bodyBytes)
	assert.Equal(t, strings.TrimSpace(res), strings.TrimSpace(bodyStr))

	logOutput := buf.String()
	assert.Contains(t, logOutput, "LOG MESSAGE • 123 • 12345")
	assert.Contains(t, logOutput, "some error")
}

