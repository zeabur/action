package environment

import (
	"slices"
	"sync"
)

// Software gets the name and version of a software installed in the current environment.
type Software interface {
	Name() string
	Version() (string, bool)
}

type softwareRegistry struct {
	software []Software
	lock     sync.RWMutex
}

var registry = softwareRegistry{
	software: nil,
	lock:     sync.RWMutex{},
}

// RegisterSoftware registers a software to the registry.
func RegisterSoftware(s Software) {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	registry.software = append(registry.software, s)
}

// GetSoftwares returns all the registered software.
func GetSoftwares() []Software {
	registry.lock.RLock()
	defer registry.lock.RUnlock()

	return slices.Clone(registry.software)
}
