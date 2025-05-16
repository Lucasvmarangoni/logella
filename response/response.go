package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Lucasvmarangoni/logella/err"
	"github.com/rs/zerolog/log"
)

type Response struct {
	Err           error      `json:"-"`
	Error         string     `json:"error"`
	Message       string     `json:"message"`
	Status        string     `json:"status"`
	RequestID     string     `json:"request_id,omitempty"`
	UserID        string     `json:"user_id,omitempty"`
	Timestamp     *time.Time `json:"timestamp,omitempty"`
	Documentation string     `json:"documentation,omitempty"`
}

func New(err error) *Response {
	e, ok := err.(*errs.Error)
	if !ok {
		return &Response{
			Status: "Internal Server Error",
			Error:  "Invalid error type passed to response.New: expected *errs.Error. Use errs.Wrap.",
		}
	}
	r := &Response{
		Err: e,
	}
	r.Status = http.StatusText(e.Code)

	r.Error = fmt.Sprintf("%v", e.ToClient())
	return r
}

func (r *Response) Log(msg string) *Response {
	parts := []string{msg}

	if r.RequestID != "" {
		parts = append(parts, r.RequestID)
	}
	if r.UserID != "" {
		parts = append(parts, r.UserID)
	}
	if errs.Unwrap(r.Err).Message != "" {
		parts = append(parts, "message:" + errs.Unwrap(r.Err).Message)
	}
	log.Error().Stack().Err(errs.Trace(r.Err).Stack()).Msgf("%s", strings.Join(parts, " • "))
	return r
}

func (r *Response) Msg(msg string) *Response {
	r.Message = msg
	return r
}

func (r *Response) Req(requestID string) *Response {
	r.RequestID = requestID
	return r
}

func (r *Response) User(userID string) *Response {
	r.UserID = userID
	return r
}

func (r *Response) Date(timestamp *time.Time) *Response {
	r.Timestamp = timestamp
	return r
}

func (r *Response) Doc(documentation string) *Response {
	r.Documentation = documentation
	return r
}

func (r *Response) Send(w http.ResponseWriter) {
	if r.Err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Unwrap(r.Err).Code)
		json.NewEncoder(w).Encode(r)
	}
}
