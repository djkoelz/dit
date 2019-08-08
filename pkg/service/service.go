package service

import (
	"net/http"
	"strconv"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	router := new(Router)
	router.mux = http.NewServeMux()

	return router
}

func (this *Router) Start(port int) {
	portString := ":" + strconv.Itoa(port)
	http.ListenAndServe(portString, this.mux)
}

func (this *Router) Register(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	this.mux.HandleFunc(pattern, handler)
}
