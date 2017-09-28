package work

import (
	"fmt"
	"sort"
	"sync"

	"github.com/heysquirrel/tribe/config"
)

type Item interface {
	GetType() string
	GetName() string
	GetDescription() string
	GetId() string
}

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

type FetchedItem struct {
	id       string
	workitem Item
	err      error
}

func (f *FetchedItem) GetId() string { return f.id }
func (f *FetchedItem) GetSummary() string {
	if f.err != nil {
		return f.err.Error()
	}

	return f.workitem.GetName()
}

func (f *FetchedItem) GetDescription() string {
	if f.err != nil {
		return f.err.Error()
	}

	return f.workitem.GetDescription()
}

type byId []*FetchedItem

func (r byId) Len() int      { return len(r) }
func (r byId) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r byId) Less(i, j int) bool {
	return r[i].id < r[j].id
}

func FetchItems(server ItemServer, ids ...string) []*FetchedItem {
	results := make([]*FetchedItem, len(ids))
	i := 0

	c := getItems(server, ids...)

	for result := range c {
		results[i], i = result, i+1
	}

	sort.Sort(sort.Reverse(byId(results)))
	return results
}

func getItems(server ItemServer, ids ...string) <-chan *FetchedItem {
	items := make(chan *FetchedItem)
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
				items <- &FetchedItem{id, workitem, err}
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
