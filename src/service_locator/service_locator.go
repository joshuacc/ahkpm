package service_locator

import "errors"

type ServiceLocator struct {
	_services map[string]interface{}
}

func NewServiceLocator() *ServiceLocator {
	return &ServiceLocator{
		_services: make(map[string]interface{}),
	}
}

var DefaultServiceLocator = NewServiceLocator()

func GetServiceLocator(maybeLocators []*ServiceLocator) *ServiceLocator {
	if len(maybeLocators) == 0 || maybeLocators[0] == nil {
		return DefaultServiceLocator
	} else {
		return maybeLocators[0]
	}
}

// Get the service with the given name.
func (sl *ServiceLocator) Get(name string) interface{} {
	service, ok := sl._services[name]
	if !ok {
		panic(errors.New("Service " + name + " not found"))
	}
	return service
}

// Add a service with the given name.
//
// If a service with the given name already exists, an error is returned.
func (sl *ServiceLocator) Add(name string, service interface{}) error {
	_, ok := sl._services[name]
	if ok {
		return errors.New("Service " + name + " already exists")
	}
	sl._services[name] = service
	return nil
}
