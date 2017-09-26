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
		return nil, err
	}

	workitem, ok := value.(WorkItem)
	if ok {
		return workitem, nil
	}

	return nil, errors.New("Unknown result")
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
