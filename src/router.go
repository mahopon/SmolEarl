package main

import (
	"net/http"
)

// Router handles all the routing for the application
type Router struct {
	controller *Controller
}

// NewRouter creates a new Router instance
func NewRouter(controller *Controller) *Router {
	return &Router{
		controller: controller,
	}
}

// Init initializes all routes
func (r *Router) Init(mux *http.ServeMux) {
	handler := LoggingMiddleware(mux)
	handler = CORSMiddleware(handler)
	mux.HandleFunc("/create", r.controller.CreateHandler)
	mux.HandleFunc("/stats/", r.controller.StatsHandler)
	mux.HandleFunc("/", r.controller.GetHandler)
}
