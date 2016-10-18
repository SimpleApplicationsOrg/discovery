package serviceDomain

import (
	"errors"
	"time"
)

func CreateDiscovery(leaseDuration time.Duration) (discovery Discovery, err error) {

	if leaseDuration == 0 {
		err = errors.New("Expire duration is mandatory")
	}

	registry := initializeRegistry(leaseDuration)

	discovery = registry

	return
}
