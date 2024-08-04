package cache

import "sync"

type CacheSet struct {
	JSCache          JSCache
	SourceIndexCache SourceIndexCache
}

func MakeCacheSet() *CacheSet {
	return &CacheSet{}
}

type SourceIndexCache struct {
	mutex           sync.Mutex
	nextSourceIndex uint32
}

func (c *SourceIndexCache) Get() uint32 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	sourceIndex := c.nextSourceIndex
	c.nextSourceIndex++
	return sourceIndex
}
