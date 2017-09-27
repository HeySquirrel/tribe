package work

import (
	"fmt"
	"github.com/heysquirrel/tribe/config"
	"sync"
)

type Item interface {
	GetType() string
	GetName() string
	GetDescription() string
	GetId() string
}

type Items []Item
type NullItem string

func (s NullItem) GetType() string        { return "" }
func (s NullItem) GetName() string        { return "" }
func (s NullItem) GetDescription() string { return "" }
func (s NullItem) GetId() string          { return string(s) }

type ItemServer interface {
	GetItem(id string) (Item, error)
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
	workitem Item
	err      error
}

func NewItemServer() (ItemServer, error) {
	servernames := config.ItemServers()
	servers := make([]ItemServer, 0)

	for _, name := range servernames {
		serverconfig := config.ItemServer(name)
		var server ItemServer
		var err error

		switch serverconfig["type"] {
		case "rally":
			server, err = NewRallyFromConfig(string(name))
			if err != nil {
				return nil, err
			}
		case "jira":
			server, err = NewJiraFromConfig(string(name))
			if err != nil {
				return nil, err
			}
		}

		servers = append(servers, NewCachingServer(server))

	}

	return NewReplicaItemServer(servers...), nil
}

func GetItems(server ItemServer, ids ...string) (Items, error) {
	results := make([]Item, 0)
	c := fetchItems(server, ids...)

	for result := range c {
		if result.err != nil {
			return results, result.err
		}
		results = append(results, result.workitem)
	}

	return results, nil
}

func fetchItems(server ItemServer, ids ...string) <-chan result {
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
				workitem, err := server.GetItem(id)
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
