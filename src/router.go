package main

import (
	"net/http"
)

// Router handles all the routing for the application

type RouterInit interface {
	Init() *http.ServeMux
}

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
func (r *Router) Init() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /stats/{id}", r.controller.StatsHandler)
	mux.HandleFunc("GET /status", r.controller.StatusHandler)
	return mux
}

type LinkRouter struct {
	controller *Controller
}

func NewLinkRouter(controller *Controller) *LinkRouter {
	return &LinkRouter{
		controller: controller,
	}
}

func (r *LinkRouter) Init() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /create", r.controller.CreateHandler)
	mux.HandleFunc("GET /{path}", r.controller.GetHandler)
	return mux
}
