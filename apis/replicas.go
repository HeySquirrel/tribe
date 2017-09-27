package apis

import (
	"time"
)

type replicas struct {
	replicas []WorkItemServer
}

func NewReplicaWorkItemServer(servers ...WorkItemServer) *replicas {
	return &replicas{servers}
}

func (m *replicas) GetWorkItem(id string) (WorkItem, error) {
	c := make(chan WorkItem)
	timeout := time.After(5 * time.Second)

	serverReplica := func(i int, id string) {
		item, err := m.replicas[i].GetWorkItem(id)
		if err == nil {
			_, ok := item.(NullWorkItem)
			if !ok {
				c <- item
			}
		}
	}

	for i := range m.replicas {
		go serverReplica(i, id)
	}

	select {
	case item := <-c:
		return item, nil
	case <-timeout:
		return NullWorkItem(id), ItemNotFoundError(id)
	}
}
