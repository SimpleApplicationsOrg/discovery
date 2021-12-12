package domain

import (
	"testing"
	"time"
)

func TestCreateDiscovery(t *testing.T) {

	if _, err := CreateDiscovery(time.Second); err != nil {
		t.Errorf("CreateRegistry() error")
	}
}

func TestRegister(t *testing.T) {

	discovery, _ := CreateDiscovery(time.Second)
	serviceName := "Service Name"
	serviceAddress := "Service address"

	if err := discovery.Register(serviceName, serviceAddress); err != nil {
		t.Errorf("Resgister failed for %s", err.Error())
	}

}

func TestRegister_Fail(t *testing.T) {

	discovery, _ := CreateDiscovery(time.Second)
	serviceName := ""
	serviceAddress := ""

	if err := discovery.Register(serviceName, serviceAddress); err == nil {
		t.Errorf("Error was expected")
	}

}

func TestFetchRegistry(t *testing.T) {

	discovery, _ := CreateDiscovery(time.Second)
	serviceName := "Service_Name"
	serviceAddress := "Service_Address"
	discovery.Register(serviceName, serviceAddress)

	if _, err := discovery.Fetch(serviceName); err != nil {
		t.Errorf("FetchRegistry error")
	}

}

func TestFetchRegistry_Fail(t *testing.T) {

	discovery, _ := CreateDiscovery(time.Second)
	serviceName := "Service Name"

	if _, err := discovery.Fetch(serviceName); err == nil {
		t.Errorf("FetchRegistry error")
	}

}

func TestUnregister(t *testing.T) {

	discovery, _ := CreateDiscovery(time.Second)
	serviceName := "Service Name"
	serviceAddress := "Service address"
	discovery.Register(serviceName, serviceAddress)
	time.Sleep(2 * time.Second)

	if _, err := discovery.Fetch(serviceName); err == nil {
		t.Errorf("Unregister failed to remove %s", serviceName)
	}

}

func TestUnregister2(t *testing.T) {

	discovery, _ := CreateDiscovery(time.Second)
	serviceName := "Service Name"
	serviceAddress := "Service address"
	discovery.Register(serviceName, serviceAddress)
	time.Sleep(2 * time.Second)

	if err := discovery.Unregister(serviceName, serviceAddress); err == nil {
		t.Errorf("Unregister failed to remove %s", serviceName)
	}

}

func Test_Given1secLease_WhenServiceLeaseExpires_ThenServiceFetchReturnsError(t *testing.T) {

	discovery, _ := CreateDiscovery(1 * time.Second)
	discovery.Register("TEST", "ADDRESS")

	time.Sleep(4 * time.Second)

	if _, err := discovery.Fetch("TEST"); err == nil {
		t.Errorf("Error expected but none received")
	}

}

func Test_Given1ServiceWith4Instances_WhenRegisterAll_ThenItCounts1ServiceWith4Instances(t *testing.T) {
	discovery, _ := CreateDiscovery(2 * time.Second)

	discovery.Register("TEST", "ADDRESS1")
	discovery.Register("TEST", "ADDRESS2")
	discovery.Register("TEST", "ADDRESS3")
	discovery.Register("TEST", "ADDRESS4")

	servicesMap := discovery.FetchAll()
	var instancesCount int
	for _, services := range servicesMap {
		instancesCount = instancesCount + len(services)
	}

	if len(servicesMap) != 1 && instancesCount != 4 {
		t.Errorf("Expected 1 service and 4 instances, but obtained %d services and %d instances", len(servicesMap), instancesCount)
	}
}

func Test_Given4ServicesWith1InstanceEach_WhenRegisterAll_ThenItCounts4ServicesAnd4Instances(t *testing.T) {
	discovery, _ := CreateDiscovery(2 * time.Second)

	discovery.Register("TEST1", "ADDRESS")
	discovery.Register("TEST2", "ADDRESS")
	discovery.Register("TEST3", "ADDRESS")
	discovery.Register("TEST4", "ADDRESS")

	servicesMap := discovery.FetchAll()
	var instancesCount int
	for _, services := range servicesMap {
		instancesCount = instancesCount + len(services)
	}

	if len(servicesMap) != 4 && instancesCount != 4 {
		t.Errorf("Expected 4 services and 4 instances, but obtained %d services and %d instances", len(servicesMap), instancesCount)
	}
}
