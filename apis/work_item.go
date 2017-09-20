package apis

import (
	"errors"
	"github.com/bluele/gcache"
	"sync"
)

type WorkItem interface {
	GetType() string
	GetName() string
	GetDescription() string
	GetId() string
}

type WorkItemServer interface {
	GetWorkItem(id string) (WorkItem, error)
}

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

type result struct {
	workitem WorkItem
	err      error
}

func GetWorkItems(server WorkItemServer, ids ...string) ([]WorkItem, error) {
	results := make([]WorkItem, 0)
	c := fetchWorkItems(server, ids...)

	for result := range c {
		if result.err != nil {
			return results, result.err
		}
		results = append(results, result.workitem)
	}

	return results, nil
}

func fetchWorkItems(server WorkItemServer, ids ...string) <-chan result {
	items := make(chan result)
	remaining := make(chan string)

	go func() {
		defer close(remaining)
		for _, id := range ids {
			remaining <- id
		}
	}()

	var wg sync.WaitGroup
	const numFetchers = 5
	wg.Add(numFetchers)

	for i := 0; i < numFetchers; i++ {
		go func() {
			for id := range remaining {
				workitem, err := server.GetWorkItem(id)
				items <- result{workitem, err}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(items)
	}()

	return items
}
