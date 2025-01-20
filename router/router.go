package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

type HTTPMethod string

const (
	POST   HTTPMethod = "POST"
	GET    HTTPMethod = "GET"
	PUT    HTTPMethod = "PUT"
	PATCH  HTTPMethod = "PATCH"
	DELETE HTTPMethod = "DELETE"
)

type Router struct {
	method HTTPMethod
	prefix string
}

func (ro *Router) InitializeRoute(r chi.Router, path string, handler http.HandlerFunc) {

	switch ro.method {
	case POST:
		r.Post(path, handler)
	case GET:
		r.Get(path, handler)
	case PUT:
		r.Put(path, handler)
	case PATCH:
		r.Patch(path, handler)
	case DELETE:
		r.Delete(path, handler)
	}
	log.Info().Str("context", "Router").Msgf("Mapped - Initialized: (%s) %s%s ", ro.method, ro.prefix, path)
}

func (ro *Router) Method(m HTTPMethod) *Router {
	ro.method = m
	return ro
}

func (ro *Router) Prefix(p string) *Router {
	ro.prefix = p
	return ro
}

func NewRouter() *Router {
	ro := &Router{}
	return ro
}
