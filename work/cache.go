package work

import (
	"errors"
	"github.com/bluele/gcache"
)

type cache struct {
	server ItemServer
	cache  gcache.Cache
}

func (c *cache) GetItem(id string) (Item, error) {
	value, err := c.cache.Get(id)
	if err != nil {
		if IsItemNotFoundError(err) {
			return NullItem(id), nil
		} else {
			return NullItem(id), err
		}
	}

	return value.(Item), nil
}

func NewCachingServer(server ItemServer) ItemServer {
	gc := gcache.New(100).
		LRU().
		LoaderFunc(func(key interface{}) (interface{}, error) {
			id, ok := key.(string)
			if ok {
				return server.GetItem(id)
			}
			return nil, errors.New("Unknown key")
		}).
		Build()

	return &cache{server, gc}
}
