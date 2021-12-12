package main

import (
	"github.com/simpleapplications/discovery/domain"
	"github.com/simpleapplications/discovery/handler"
	"github.com/simpleapplications/discovery/router"
	"net/http"
	"time"
)

const ServiceHostKeepTime = 5 * time.Second

func main() {
	discovery, err := domain.CreateDiscovery(ServiceHostKeepTime)
	if err != nil {
		panic(err)
	}
	h := handler.Discovery{Registry: discovery}

	r := router.New()
	r.Add(http.MethodPost, "/services", h.Register)
	r.Add(http.MethodGet, "/services", h.FetchAll)
	r.Add(http.MethodGet, "/services/{id}", h.Fetch)
	r.Add(http.MethodDelete, "/services/{id}/hosts/{url}", h.Unregister)

	r.Start()
}
