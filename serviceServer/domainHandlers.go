package serviceServer

import (
	"encoding/json"
	"fmt"
	"github.com/SimpleApplicationsOrg/discovery/serviceDomain"
	"net/http"
	"strings"
)

type ResponseMessage struct {
	Message string `json:"response"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	parameters := findPathParameters("register", r)

	var serviceName string
	if len(parameters) > 0 {
		serviceName = parameters[0]
	}

	var serviceAddress string
	if len(parameters) > 1 {
		serviceAddress = parameters[1]
	}

	discovery := r.Context().Value("discovery").(serviceDomain.Discovery)

	if err := discovery.Register(serviceName, serviceAddress); err != nil {
		http.Error(w, fmt.Sprintf("Error registering service: %s", err.Error()), http.StatusNotFound)

	} else {
		response := ResponseMessage{Message: fmt.Sprintf("Service %s (%s) registered", serviceName, serviceAddress)}
		json.NewEncoder(w).Encode(response)

	}

	return
}

type ResponseFetch struct {
	ServiceName string `json:"service_name"`
	ServiceAddress string `json:"service_address"`
}

func FetchHandler(w http.ResponseWriter, r *http.Request) {

	parameters := findPathParameters("fetch", r)

	var serviceName string
	if len(parameters) > 0 {
		serviceName = parameters[0]
	}

	discovery := r.Context().Value("discovery").(serviceDomain.Discovery)

	if service, err := discovery.Fetch(serviceName); err != nil {
		http.Error(w, fmt.Sprintf("Error fetching service %s: %s", serviceName, err.Error()), http.StatusNotFound)
	} else {
		response := ResponseFetch{ServiceName:service.Name, ServiceAddress:service.Address}
		json.NewEncoder(w).Encode(response)
	}

	return
}

func RenewHandler(w http.ResponseWriter, r *http.Request) {

	parameters := findPathParameters("renew", r)

	var serviceName string
	if len(parameters) > 0 {
		serviceName = parameters[0]
	}

	var serviceAddress string
	if len(parameters) > 1 {
		serviceAddress = parameters[1]
	}

	discovery := r.Context().Value("discovery").(serviceDomain.Discovery)

	if err := discovery.Renew(serviceName, serviceAddress); err != nil {
		http.Error(w, fmt.Sprintf("Error renewing service: %s", err.Error()), http.StatusNotFound)

	} else {
		response := ResponseMessage{Message: fmt.Sprintf("Service %s (%s) renewed", serviceName, serviceAddress)}
		json.NewEncoder(w).Encode(response)

	}

	return
}

type ResponseFetchAll struct {
	Instances map[string]serviceDomain.ServiceInstances `json:"instances"`
}

func FetchAllHandler(w http.ResponseWriter, r *http.Request) {

	discovery := r.Context().Value("discovery").(serviceDomain.Discovery)

	instances := discovery.FetchAll()

	response := ResponseFetchAll{Instances: instances}
	json.NewEncoder(w).Encode(response)

	return
}

func UnregisterHandler(w http.ResponseWriter, r *http.Request) {

	parameters := findPathParameters("unregister", r)

	var serviceName string
	if len(parameters) > 0 {
		serviceName = parameters[0]
	}

	var serviceAddress string
	if len(parameters) > 1 {
		serviceAddress = parameters[1]
	}

	discovery := r.Context().Value("discovery").(serviceDomain.Discovery)

	if err := discovery.Unregister(serviceName, serviceAddress); err != nil {
		http.Error(w, fmt.Sprintf("Error unregistering service: %s", err.Error()), http.StatusNotFound)

	} else {
		response := ResponseMessage{Message: fmt.Sprintf("Service %s (%s) unregistered", serviceName, serviceAddress)}
		json.NewEncoder(w).Encode(response)

	}

	return
}

func findPathParameters(parameter string, r *http.Request) []string {
	pathList := strings.Split(r.URL.Path, "/")

	for i, word := range pathList {
		if word == parameter {
			return pathList[i+1:]
		}
	}
	return nil
}
