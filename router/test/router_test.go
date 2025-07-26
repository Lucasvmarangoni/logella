package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Lucasvmarangoni/logella/router"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
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

type UsersRouter struct {
	router *router.Router
	prefix string
}

func NewUsersRouter(
	router *router.Router,
) *UsersRouter {
	return &UsersRouter{
		router: router,
		prefix: router.Prefix,
	}
}

func (u *UsersRouter) testRouter() {
	R := u.router
	R.Route("/user", func() {
		R.Get("/get", handler_user)
		R.Post("/post", handler_user1)

		R.Group(func() {
			R.Use(httprate.Limit(
				4,
				60*time.Minute,
				httprate.WithKeyFuncs(httprate.KeyByRealIP, httprate.KeyByEndpoint),
				httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				}),
			))
			R.Post("/put", handler_user)
		})
	})
}

func (u *UsersRouter) testRouter1() {
	R := u.router
	R.Route("/user1", func() {
		R.Get("/get1", handler_user1)
		R.Post("/post1", handler_user)
		R.Group(func() {
			R.Use(httprate.Limit(
				4,
				60*time.Minute,
				httprate.WithKeyFuncs(httprate.KeyByRealIP, httprate.KeyByEndpoint),
				httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				}),
			))
			R.Post("/put1", handler_group)
		})
	})
}

func TestRoutes(t *testing.T) {

	log.Logger = zerolog.New(&buf).With().Timestamp().Logger()

	tests := []struct {
		method string
		url    string
	}{
		{"GET", "/user/get"},
		{"POST", "/user/post"},
		{"POST", "/user/put"},
		{"GET", "/user1/get1"},
		{"POST", "/user1/post1"},
		{"POST", "/user1/put1"},
	}

	R := router.NewRouter()
	router := NewUsersRouter(R)

	R.Use(middleware.Logger, middleware.Recoverer)
	router.testRouter()
	R.Group(func() {
		router.testRouter1()
	})

	for _, test := range tests {
		t.Run(test.method+" "+test.url, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.url, nil)
			r := httptest.NewRecorder()
			R.Chi.ServeHTTP(r, req)

			expectedLog := "Mapped - Initialized: " + "(" + test.method + ")" + " " + test.url
			assert.Contains(t, buf.String(), expectedLog)
		})
	}
}
