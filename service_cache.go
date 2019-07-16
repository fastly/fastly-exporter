package main

import "sync/atomic"

// serviceCache is a simple wrapper around a map of service ID to name and
// current version.
type serviceCache struct {
	v atomic.Value
}

type nameVersion struct {
	name    string
	version string
}

// newServiceCache returns an empty and usable service cache.
func newServiceCache() *serviceCache {
	var v atomic.Value
	v.Store(map[string]nameVersion{})
	return &serviceCache{v}
}

// update the complete mapping of service IDs to names and versions.
func (c *serviceCache) update(services map[string]nameVersion) {
	c.v.Store(services)
}

// resolve a service ID to its corresponding name and version. If the service ID
// isn't found, the ID itself is returned as the name, and "unknown" is returned
// for the version.
func (c *serviceCache) resolve(id string) (name, version string) {
	if nv, ok := c.v.Load().(map[string]nameVersion)[id]; ok {
		return nv.name, nv.version
	}
	return id, "unknown"
}
