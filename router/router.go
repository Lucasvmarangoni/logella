package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

type HTTPMethod string

const (
	post   HTTPMethod = "POST"
	get    HTTPMethod = "GET"
	put    HTTPMethod = "PUT"
	patch  HTTPMethod = "PATCH"
	delete HTTPMethod = "DELETE"
)

type Router struct {
	Chi    chi.Router
	mux    *chi.Mux
	Prefix string
}

func (r *Router) Route(pattern string, fn func()) chi.Router {
	if fn == nil {
		panic(fmt.Sprintf("chi: attempting to Route() a nil subrouter on '%s'", pattern))
	}
	r.Prefix = pattern
	subRouter := NewRouter()
	r.Chi.Mount(pattern, subRouter.Chi)
	fn()
	return subRouter.Chi
}

func (r *Router) Group(fn func()) *Router {
	if fn != nil {
		r.mux.Group(func(chi.Router) {
			fn()
		})

	}
	return r
}

func (r *Router) Use(ms ...func(http.Handler) http.Handler) *Router {
	r.mux.Use(ms...)
	return r
}

func (r *Router) Post(pattern string, handler http.HandlerFunc) {
	r.Chi.Post(pattern, handler)
	log.Info().Str("context", "Router").Msgf("Mapped - Initialized: (%s) %s%s ", post, r.Prefix, pattern)
}

func (r *Router) Get(pattern string, handler http.HandlerFunc) {
	r.Chi.Get(pattern, handler)
	log.Info().Str("context", "Router").Msgf("Mapped - Initialized: (%s) %s%s ", get, r.Prefix, pattern)
}

func (r *Router) Put(pattern string, handler http.HandlerFunc) {
	r.Chi.Put(pattern, handler)
	log.Info().Str("context", "Router").Msgf("Mapped - Initialized: (%s) %s%s ", put, r.Prefix, pattern)
}

func (r *Router) Path(pattern string, handler http.HandlerFunc) {
	r.Chi.Patch(pattern, handler)
	log.Info().Str("context", "Router").Msgf("Mapped - Initialized: (%s) %s%s ", patch, r.Prefix, pattern)
}

func (r *Router) Delete(pattern string, handler http.HandlerFunc) {
	r.Chi.Delete(pattern, handler)
	log.Info().Str("context", "Router").Msgf("Mapped - Initialized: (%s) %s%s ", delete, r.Prefix, pattern)
}

func NewRouter() *Router {
	r := &Router{
		Chi: chi.NewRouter(),
		mux: chi.NewMux(),
	}
	return r
}
