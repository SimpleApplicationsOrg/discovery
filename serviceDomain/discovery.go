package domain

type Discovery interface {
	Register(string, string) error
	Renew(string, string) error
	Fetch(string) (Service, error)
	FetchAll() map[string]ServiceInstances
	Unregister(string, string) error
	Close()
}
