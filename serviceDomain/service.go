package domain

import "time"

type ServiceInstances struct {
	Instances map[string]Service
}

type Service struct {
	Name         string
	Address      string
	LastRegister time.Time
}
