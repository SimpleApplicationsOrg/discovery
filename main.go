package main

import (
	"context"
	"github.com/SimpleApplicationsOrg/discovery/serviceDomain"
	"github.com/SimpleApplicationsOrg/discovery/serviceServer"
	"log"
	"time"
)

func main() {
	server := serviceServer.NewSimpleServer()

	discovery, err := serviceDomain.CreateDiscovery(30 * time.Second)
	if err != nil {
		panic("Error creating discovery registry: " + err.Error())
	}
	log.Println("Starting registry...")
	defer discovery.Close()

	discovery_ctx := context.WithValue(context.Background(), "discovery", discovery)

	server.AddMappingWithContext("/register/", "PUT", serviceServer.RegisterHandler, discovery_ctx)
	server.AddMappingWithContext("/fetch/", "GET", serviceServer.FetchHandler, discovery_ctx)
	server.AddMappingWithContext("/renew/", "PUT", serviceServer.RenewHandler, discovery_ctx)
	server.AddMappingWithContext("/fetchAll", "GET", serviceServer.FetchAllHandler, discovery_ctx)
	server.AddMappingWithContext("/unregister/", "PUT", serviceServer.UnregisterHandler, discovery_ctx)

	server.Start()
}
