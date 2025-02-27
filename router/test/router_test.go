package router_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Lucasvmarangoni/logella/router"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

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
	R.Route("/user", func(R *router.Router) {
		R.Get("/get", handler_user)
		R.Post("/post", handler_user1)
		R.Post("/put", handler_user)
	})

	R.Route("/user1", func(R *router.Router) {
		R.Get("/get1", handler_user1)
		R.Post("/post1", handler_user)

		R.Chi.Group(func(r chi.Router) {
			r.Post("/put1", handler_group)
		})
	})
	return R.Chi
}

func TestRoutes(t *testing.T) {
	router := testRouter()
	tests := []struct {
		method string
		url    string
		status int
		body   string
	}{
		{"GET", "/user/get", http.StatusOK, "\"handler user 1\""},
		{"POST", "/user/post", http.StatusOK, "\"handler user 1\""},
		{"POST", "/user/put", http.StatusOK, "\"handler user 1\""},
		{"GET", "/user1/get1", http.StatusOK, "\"handler user 1\""},
		{"POST", "/user1/post1", http.StatusOK, "\"handler user 1\""},
		{"POST", "/user1/put1", http.StatusOK, "\"handler group\""},
	}

	for _, test := range tests {
		t.Run(test.method+" "+test.url, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.url, nil)
			r := httptest.NewRecorder()
			router.ServeHTTP(r, req)

			assert.Equal(t, test.status, r.Code, "%s %s: expected status %d, got %d", test.method, test.url, test.status, r.Code)

			respBody, _ := io.ReadAll(r.Body)
			assert.Equal(t, test.body+"\n", string(respBody), "%s %s: expected body %s, got %s", test.method, test.url, test.body, string(respBody))
		})
	}
}
