package main

import "sync/atomic"

// nameCache is a simple wrapper around a map of service ID to service name.
type nameCache struct {
	v atomic.Value
}

// newNameCache returns an empty and usable name cache.
func newNameCache() *nameCache {
	var v atomic.Value
	v.Store(map[string]string{})
	return &nameCache{v}
}

// update the complete mapping of service IDs to names.
func (c *nameCache) update(names map[string]string) {
	c.v.Store(names)
}

// resolve a service ID to its corresponding name.
// If the service ID isn't found, the ID itself is returned as the name.
func (c *nameCache) resolve(id string) (name string) {
	if name, ok := c.v.Load().(map[string]string)[id]; ok {
		return name
	}
	return id
}
