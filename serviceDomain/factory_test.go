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
