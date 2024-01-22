package router

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

type Router struct {
	method        string	
	prefix        string
}

func (ro *Router) InitializeRoute(r chi.Router, path string, handler http.HandlerFunc) {

	switch strings.ToUpper(ro.method) {
	case "POST":
		r.Post(path, handler)
	case "GET":
		r.Get(path, handler)
	case "PUT":
		r.Put(path, handler)
	case "PATCH":
		r.Patch(path, handler)
	case "DELETE":
		r.Delete(path, handler)
	}
	log.Info().Str("context", "Router").Msgf("Mapped - Initialized: (%s) %s%s ", ro.method, ro.prefix, path)
}

func (ro *Router) Method(m string) *Router {
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