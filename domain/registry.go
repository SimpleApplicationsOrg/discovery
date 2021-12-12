package domain

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Service struct {
	Name         string
	Address      string
	LastRegister time.Time
}

type ServiceRegistry struct {
	mu            sync.Mutex
	services      map[string][]Service
	leaseDuration time.Duration
	closeDaemon   chan bool
}

func CreateDiscovery(leaseDuration time.Duration) (*ServiceRegistry, error) {
	if leaseDuration == 0 {
		return nil, errors.New("expire duration is mandatory")
	}
	return initializeRegistry(leaseDuration), nil
}

func initializeRegistry(lease time.Duration) *ServiceRegistry {
	registry := ServiceRegistry{services: make(map[string][]Service), leaseDuration: lease}
	registry.closeDaemon = make(chan bool)

	go registry.unregisterDaemon()
	return &registry
}

func (r *ServiceRegistry) Close() {
	r.closeDaemon <- true
}

func (r *ServiceRegistry) Register(name string, address string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(name) == 0 || len(address) == 0 {
		return errors.New("name and Address are mandatory")
	}

	services, ok := r.services[name]
	if !ok {
		r.services[name] = make([]Service, 0)
	}

	for _, svc := range services {
		if svc.Address == address {
			log.Println("Updating address " + address + " for service " + name)
			svc.LastRegister = time.Now()

			r.services[name] = services
			return nil
		}
	}

	log.Println("Registering address " + address + " for service " + name)
	instance := Service{Name: name, Address: address, LastRegister: time.Now()}
	r.services[name] = append(r.services[name], instance)

	return nil
}

func (r *ServiceRegistry) Fetch(name string) (Service, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(name) == 0 {
		return Service{}, errors.New("name is mandatory")
	}

	services, ok := r.services[name]
	if !ok {
		return Service{}, errors.New("service not found")
	}

	return services[rand.Intn(len(services))], nil
}

func (r *ServiceRegistry) FetchAll() map[string][]Service {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Println("Getting all instances")
	return r.services
}

func (r *ServiceRegistry) Unregister(name string, address string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Println("unregister", name, address)

	if len(name) == 0 {
		return errors.New("name is mandatory")
	}

	services, ok := r.services[name]
	if !ok {
		return errors.New("service not found")
	}

	for i, svc := range services {
		if svc.Address == address {
			log.Println("Removing instance " + address)
			r.services[name] = remove(r.services[name], i)
			break
		}
	}

	if len(r.services[name]) == 0 {
		log.Println("Removing service " + name)
		delete(r.services, name)
	}

	return nil
}

func (r *ServiceRegistry) unregisterDaemon() {
	for {
		select {
		case closeDaemon := <-r.closeDaemon:
			if closeDaemon {
				return
			}
		default:
			time.Sleep(r.leaseDuration)
			r.unregisterExpiredServices()
		}
	}
}

func (r *ServiceRegistry) unregisterExpiredServices() {
	serviceMap := r.FetchAll()

	for _, services := range serviceMap {
		for _, instance := range services {
			expirationTime := instance.LastRegister.Add(r.leaseDuration)
			if time.Now().After(expirationTime) {
				_ = r.Unregister(instance.Name, instance.Address)
			}
		}
	}
}

func remove(s []Service, i int) []Service {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
