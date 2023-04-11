package router

import "net/http"

type Router struct {
	Handler http.Handler
	Mux     http.ServeMux
}

func NewRouter(handeler http.Handler, mux http.ServeMux) *Router {
	return &Router{
		Handler: handeler,
		Mux:     mux,
	}
}
