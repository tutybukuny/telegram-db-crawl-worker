package container

import (
	"log"

	"github.com/golobby/container/v3"
)

func NamedSingleton(name string, resolver any) {
	err := container.NamedSingleton(name, resolver)
	if err != nil {
		log.Fatalf("cannot bind object with name %s: %v", name, err)
	}
}

func Fill(receiver any) {
	err := container.Fill(receiver)
	if err != nil {
		log.Fatalf("cannot fill: %v", err)
	}
}

func Resolve(abstraction any) {
	err := container.Resolve(abstraction)
	if err != nil {
		log.Fatalf("cannot resolve: %v", err)
	}
}

func NamedResolve(abstraction any, name string) {
	err := container.NamedResolve(abstraction, name)
	if err != nil {
		log.Fatalf("cannot named resolve with name %s: %v", name, err)
	}
}
