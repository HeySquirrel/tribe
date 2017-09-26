package apis

import (
	"fmt"
	"sync"
)

type WorkItem interface {
	GetType() string
	GetName() string
	GetDescription() string
	GetId() string
}

type NotFoundWorkItem string

func (s NotFoundWorkItem) GetType() string        { return "" }
func (s NotFoundWorkItem) GetName() string        { return "" }
func (s NotFoundWorkItem) GetDescription() string { return "" }
func (s NotFoundWorkItem) GetId() string          { return string(s) }

type ItemNotFoundError string

func (s ItemNotFoundError) Error() string {
	return fmt.Sprintf("'%s' was not found.", string(s))
}

type WorkItemServer interface {
	GetWorkItem(id string) (WorkItem, error)
}

type WorkItems []WorkItem

type result struct {
	workitem WorkItem
	err      error
}

func GetWorkItems(server WorkItemServer, ids ...string) (WorkItems, error) {
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
				if err != nil && err == ItemNotFoundError(id) {
					items <- result{NotFoundWorkItem(id), nil}
				} else {
					items <- result{workitem, err}
				}
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
