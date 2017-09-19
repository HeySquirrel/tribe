package apis

import (
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
