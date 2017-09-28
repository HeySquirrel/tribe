package work

import (
	"time"
)

type replicas struct {
	replicas []ItemServer
}

func NewReplicaItemServer(servers ...ItemServer) *replicas {
	return &replicas{servers}
}

func (m *replicas) GetItem(id string) (Item, error) {
	c := make(chan Item)
	timeout := time.After(5 * time.Second)

	serverReplica := func(i int, id string) {
		item, err := m.replicas[i].GetItem(id)
		if err == nil {
			c <- item
		}
	}

	for i := range m.replicas {
		go serverReplica(i, id)
	}

	select {
	case item := <-c:
		return item, nil
	case <-timeout:
		return nil, ItemNotFoundError(id)
	}
}
