package apis

import (
	"errors"
	"github.com/bluele/gcache"
)

type cache struct {
	server WorkItemServer
	cache  gcache.Cache
}

func (c *cache) GetWorkItem(id string) (WorkItem, error) {
	value, err := c.cache.Get(id)
	if err != nil {
		if IsItemNotFoundError(err) {
			return NullWorkItem(id), nil
		} else {
			return NullWorkItem(id), err
		}
	}

	return value.(WorkItem), nil
}

func NewCachingServer(server WorkItemServer) WorkItemServer {
	gc := gcache.New(100).
		LRU().
		LoaderFunc(func(key interface{}) (interface{}, error) {
			id, ok := key.(string)
			if ok {
				return server.GetWorkItem(id)
			}
			return nil, errors.New("Unknown key")
		}).
		Build()

	return &cache{server, gc}
}
