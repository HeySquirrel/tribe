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

type WorkItems []WorkItem
type NullWorkItem string

func (s NullWorkItem) GetType() string        { return "" }
func (s NullWorkItem) GetName() string        { return "" }
func (s NullWorkItem) GetDescription() string { return "" }
func (s NullWorkItem) GetId() string          { return string(s) }

type WorkItemServer interface {
	GetWorkItem(id string) (WorkItem, error)
}

type ItemNotFoundError string

func (s ItemNotFoundError) Error() string {
	return fmt.Sprintf("'%s' was not found.", string(s))
}

func IsItemNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(ItemNotFoundError)
	return ok
}

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
				if IsItemNotFoundError(err) {
					items <- result{workitem, nil}
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
