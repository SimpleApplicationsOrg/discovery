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
	services      sync.Map
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
	registry := ServiceRegistry{leaseDuration: lease}
	registry.closeDaemon = make(chan bool)

	go registry.unregisterDaemon()
	return &registry
}

func (r *ServiceRegistry) Close() {
	r.closeDaemon <- true
}

func (r *ServiceRegistry) Register(name string, address string) error {
	if len(name) == 0 || len(address) == 0 {
		return errors.New("name and Address are mandatory")
	}

	var services []Service
	servicesI, ok := r.services.Load(name)
	if !ok {
		r.services.Store(name, make([]Service, 0))
	} else {
		services = servicesI.([]Service)
	}

	for i, svc := range services {
		if svc.Address == address {
			log.Println("Updating address " + address + " for service " + name)
			services[i].LastRegister = time.Now()

			r.services.Store(name, services)
			return nil
		}
	}

	log.Println("Registering address " + address + " for service " + name)
	instance := Service{Name: name, Address: address, LastRegister: time.Now()}
	r.services.Store(name, append(services, instance))

	return nil
}

func (r *ServiceRegistry) Fetch(name string) (Service, error) {
	if len(name) == 0 {
		return Service{}, errors.New("name is mandatory")
	}

	var services []Service
	servicesI, ok := r.services.Load(name)
	if !ok {
		return Service{}, errors.New("service not found")
	} else {
		services = servicesI.([]Service)
	}

	return services[rand.Intn(len(services))], nil
}

func (r *ServiceRegistry) FetchAll() map[string][]Service {
	servicesMap := make(map[string][]Service)
	r.services.Range(func(key, value interface{}) bool {
		servicesMap[key.(string)] = value.([]Service)
		return true
	})

	log.Println("Getting all instances")
	return servicesMap
}

func (r *ServiceRegistry) Unregister(name string, address string) error {
	log.Println("unregister", name, address)

	if len(name) == 0 {
		return errors.New("name is mandatory")
	}

	var services []Service
	servicesI, ok := r.services.Load(name)
	if !ok {
		return errors.New("service not found")
	} else {
		services = servicesI.([]Service)
	}

	toRemove := -1
	for i, svc := range services {
		if svc.Address == address {
			log.Println("Removing instance " + address)
			toRemove = i
			break
		}
	}
	if toRemove < 0 {
		return errors.New("address not found")
	}

	services[toRemove] = services[len(services)-1]
	services = services[:len(services)-1]
	r.services.Store(name, services)
	
	if len(services) == 0 {
		log.Println("Removing service " + name)
		r.services.Delete(name)
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
