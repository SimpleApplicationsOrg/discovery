package handler

import (
	"encoding/json"
	"fmt"
	"github.com/simpleapplications/discovery/domain"
	"io"
	"log"
	"net/http"
)

type RegisterRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type ResponseFetch struct {
	Address string `json:"address"`
}

type ResponseFetchAll struct {
	Instances map[string][]domain.Service `json:"instances"`
}

type Discovery struct {
	Registry *domain.ServiceRegistry
}

func (d *Discovery) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := parseBody(r, &req); err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := d.Registry.Register(req.Name, req.Address); err != nil {
		http.Error(w, fmt.Sprintf("Error registering service: %s", err.Error()), http.StatusNotFound)
		return
	}

	return
}

func (d *Discovery) Fetch(w http.ResponseWriter, r *http.Request) {
	serviceName := getPathParams(r, "id")
	service, err := d.Registry.Fetch(serviceName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching service %s: %s", serviceName, err.Error()), http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(ResponseFetch{Address: service.Address})
	return
}

func (d *Discovery) FetchAll(w http.ResponseWriter, r *http.Request) {
	instances := d.Registry.FetchAll()

	response := ResponseFetchAll{Instances: instances}
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)

	return
}

func (d *Discovery) Unregister(w http.ResponseWriter, r *http.Request) {
	serviceName := getPathParams(r, "id")
	serviceAddress := getPathParams(r, "url")
	if err := d.Registry.Unregister(serviceName, serviceAddress); err != nil {
		http.Error(w, fmt.Sprintf("Error unregistering service: %s", err.Error()), http.StatusNotFound)
	}

	return
}

func parseBody(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to get body, %+v", err)
	}
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to parse body, %+v", err)
	}

	return nil
}

func getPathParams(r *http.Request, k string) string {
	params, ok := r.Context().Value("pathParams").(map[string]string)
	if !ok {
		log.Printf("failed get parameters")
		return ""
	}

	resp, ok := params[k]
	if !ok {
		log.Printf("failed get path parameter %s", k)
		return ""
	}

	return resp
}
