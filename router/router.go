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
	Prefix string
}

func (r *Router) Route(pattern string, fn func()) chi.Router {
	if fn == nil {
		panic(fmt.Sprintf("chi: attempting to Route() a nil subrouter on '%s'", pattern))
	}
	r.Prefix = pattern
	r.Chi.Route(pattern, func(sub chi.Router) {
		prev := r.Chi
		r.Chi = sub
		fn()
		r.Chi = prev
	})
	return r.Chi
}

func (r *Router) Group(fn func()) {
	prev := r.Chi
	r.Chi.Group(func(gr chi.Router) {
		r.Chi = gr
		fn()
	})
	r.Chi = prev
}

func (r *Router) Use(ms ...func(http.Handler) http.Handler) {
	r.Chi.Use(ms...)
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
	}
	return r
}
