package server

import (
	"context"
	"reflect"
)

type Descriptor struct {
	Name     string
	Instance Service
}

type Service interface {
	Init() error
}

type BackgroundService interface {
	Run(ctx context.Context) error
}

var services []*Descriptor

func RegisterService(instance Service) {
	services = append(services, &Descriptor{
		Name:     reflect.TypeOf(instance).Elem().Name(),
		Instance: instance,
	})
}

func Register(descriptor *Descriptor) {
	services = append(services, descriptor)
}

func GetServices() []*Descriptor {
	return services
}
