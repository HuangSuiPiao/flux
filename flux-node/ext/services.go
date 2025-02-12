package ext

import (
	"fmt"
	"github.com/bytepowered/flux/flux-node"
	"sync"
)

var (
	serviceNotFound flux.Service
	servicesMap     = new(sync.Map)
)

func RegisterServiceByID(id string, service flux.Service) {
	servicesMap.Store(id, service)
}

// RegisterService store transporter service
func RegisterService(service flux.Service) {
	id := _ensureServiceID(&service)
	RegisterServiceByID(id, service)
}

// Services 返回全部注册的Service
func Services() map[string]flux.Service {
	out := make(map[string]flux.Service, 512)
	endpoints.Range(func(key, value interface{}) bool {
		out[key.(string)] = value.(flux.Service)
		return true
	})
	return out
}

// ServiceByID load transporter service by serviceId
func ServiceByID(serviceID string) (flux.Service, bool) {
	v, ok := servicesMap.Load(serviceID)
	if ok {
		return v.(flux.Service), true
	}
	return serviceNotFound, false
}

// RemoveServiceByID remove transporter service by serviceId
func RemoveServiceByID(serviceID string) {
	servicesMap.Delete(serviceID)
}

// HasServiceByID check service exists by service id
func HasServiceByID(serviceID string) bool {
	_, ok := servicesMap.Load(serviceID)
	return ok
}

func _ensureServiceID(service *flux.Service) string {
	id := service.ServiceId
	if id == "" {
		id = service.Interface + ":" + service.Method
	}
	if len(id) < len("a:b") {
		panic(fmt.Sprintf("Transporter must has an Id, service: %+v", service))
	}
	return id
}
