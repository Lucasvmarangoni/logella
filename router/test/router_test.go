package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Lucasvmarangoni/logella/router"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

var buf bytes.Buffer

func handler_user(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("handler user 1")
	w.WriteHeader(http.StatusOK)
}

func handler_user1(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("handler user 1")
	w.WriteHeader(http.StatusOK)
}

func handler_group(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("handler group")
	w.WriteHeader(http.StatusOK)
}

func testRouter() chi.Router {
	R := router.NewRouter()
	R.Route("/user", func() {
		R.Get("/get", handler_user)
		R.Post("/post", handler_user1)
		R.Post("/put", handler_user)
	})

	R.Route("/user1", func() {
		R.Get("/get1", handler_user1)
		R.Post("/post1", handler_user)
		R.Group(func(){
			R.Post("/put1", handler_group)
		})
	})
	return R.Chi
}

func TestRoutes(t *testing.T) {

	log.Logger = zerolog.New(&buf).With().Timestamp().Logger()

	tests := []struct {
		method  string
		url     string
	}{
		{"GET", "/user/get"},
		{"POST", "/user/post"},
		{"POST", "/user/put"},
		{"GET", "/user1/get1"},
		{"POST", "/user1/post1"},
		{"POST", "/user1/put1",},
	}

	router := testRouter()
	for _, test := range tests {
		t.Run(test.method+" "+test.url, func(t *testing.T) {			
			req := httptest.NewRequest(test.method, test.url, nil)
			r := httptest.NewRecorder()
			router.ServeHTTP(r, req)	

			expectedLog := "Mapped - Initialized: " + "(" + test.method + ")" + " " + test.url
			assert.Contains(t, buf.String(), expectedLog)
		})
	}
}
