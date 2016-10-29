package serviceDomain

import (
	"errors"
	"log"
	"math/rand"
	"time"
	"sync"
)

var mutex sync.Mutex

type serviceRegistry struct {
	services      map[string]ServiceInstances
	leaseDuration time.Duration
	closeDaemon   chan bool
}

func initializeRegistry(lease time.Duration) serviceRegistry {
	registry := serviceRegistry{services: make(map[string]ServiceInstances), leaseDuration: lease}
	registry.closeDaemon = make(chan bool)

	go registry.unregisterDaemon()
	return registry
}

func (registry serviceRegistry) Close() {
	registry.closeDaemon <- true
}

func (registry serviceRegistry) Register(name string, address string) (err error) {


	if len(name) == 0 || len(address) == 0 {
		err = errors.New("Name and Address are mandatory")
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	service, ok := registry.services[name]
	if !ok {
		service = ServiceInstances{}
		registry.services[name] = service
	}

	if service.Instances == nil {
		service.Instances = make(map[string]Service)
	}

	instance := Service{Name: name, Address: address, LastRegister: time.Now()}
	service.Instances[address] = instance
	registry.services[name] = service

	return
}

func (registry serviceRegistry) Renew(name string, address string) (err error) {

	if len(name) == 0 {
		err = errors.New("Name is mandatory")
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	service, ok := registry.services[name]

	if !ok {
		err = errors.New("Service not found")
		return

	}

	instance, ok := service.Instances[address]

	instance.LastRegister = time.Now()

	return

}

func (registry serviceRegistry) Fetch(name string) (instance Service, err error) {

	if len(name) == 0 {
		err = errors.New("Name is mandatory")
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	service, ok := registry.services[name]

	if !ok || len(service.Instances) == 0 {
		err = errors.New("Service not found")
		return
	}

	loadBalanceResolution := rand.Intn(len(service.Instances))
	index := 0
	for _, service := range service.Instances {
		if index == loadBalanceResolution {
			log.Println("Getting address " + service.Address + " for service " + service.Name)
			instance = service
		}
		index = index + 1
	}

	return
}

func (registry serviceRegistry) FetchAll() map[string]ServiceInstances {
	return registry.services
}

func (registry serviceRegistry) Unregister(name string, address string) (err error) {


	if len(name) == 0 {
		err = errors.New("Name is mandatory")
		return
	}

	service, ok := registry.services[name]

	if ok {
		instance, ok := service.Instances[address]

		if ok {
			log.Println("Removing instance " + instance.Address)
			delete(service.Instances, instance.Address)
		}

		if len(service.Instances) == 0 {
			log.Println("Removing service " + instance.Name)
			delete(registry.services, instance.Name)
		}

	} else {
		err = errors.New("Service not found")
	}

	return

}

func (registry serviceRegistry) unregisterDaemon() {

	for {
		select {
		case closeDaemon := <-registry.closeDaemon:
			if closeDaemon {
				return
			}
		default:
			time.Sleep(registry.leaseDuration)
			registry.unregisterExpiredServices()

		}
	}

}

func (registry serviceRegistry) unregisterExpiredServices() {

	mutex.Lock()
	defer mutex.Unlock()

	services := make(map[string]ServiceInstances)
	for k, v := range registry.services {

		services[k] = v
	}

	for _, service := range services {

		for _, instance := range service.Instances {

			expirationTime := instance.LastRegister.Add(registry.leaseDuration)

			if time.Now().After(expirationTime) {
				registry.Unregister(instance.Name, instance.Address)
			}

		}

	}

}
