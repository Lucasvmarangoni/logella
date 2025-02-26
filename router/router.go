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
	Chi    chi.Router
	method HTTPMethod
	prefix string
}

func (r *Router) Map(pattern string, handler http.HandlerFunc) {

	switch r.method {
	case POST:
		r.Chi.Post(pattern, handler)
	case GET:
		r.Chi.Get(pattern, handler)
	case PUT:
		r.Chi.Put(pattern, handler)
	case PATCH:
		r.Chi.Patch(pattern, handler)
	case DELETE:
		r.Chi.Delete(pattern, handler)
	}
	log.Info().Str("context", "Router").Msgf("Mapped - Initialized: (%s) %s%s ", r.method, r.prefix, pattern)
}

func (r *Router) Method(m HTTPMethod) *Router {
	r.method = m
	return r
}

func (r *Router) Prefix(p string) *Router {
	r.prefix = p
	return r
}

func NewRouter() *Router {
	r := &Router{
		Chi: chi.NewRouter(),
	}
	return r
}
