package main

import "sync/atomic"

type nameCache struct {
	v atomic.Value
}

func newNameCache() *nameCache {
	var v atomic.Value
	v.Store(map[string]string{})
	return &nameCache{v}
}

func (c *nameCache) update(names map[string]string) {
	c.v.Store(names)
}

func (c *nameCache) resolve(id string) (name string) {
	if name, ok := c.v.Load().(map[string]string)[id]; ok {
		return name
	}
	return id
}
