package apis

import (
	"testing"
)

func setupRally(t *testing.T) WorkItemServer {
	api, err := NewRallyFromConfig("rally1")
	if err != nil {
		t.Fatal(err)
	}
	return api
}

func TestRallyGetWorkItem(t *testing.T) {
	api := setupRally(t)
	expectedType := "HierarchicalRequirement"

	item, err := api.GetWorkItem("S113541")
	if err != nil {
		t.Fatal(err)
	}

	if expectedType != item.GetType() {
		t.Fatalf("Expected workitem to be of type '%s', but was: '%s'", expectedType, item.GetType())
	}
}

func TestRallyNotFoundWorkItem(t *testing.T) {
	api := setupRally(t)
	itemid := "NOTID"

	item, err := api.GetWorkItem(itemid)
	if err == nil || err != ItemNotFoundError(itemid) {
		t.Fatal("Expected ItemNotFoundError")
	}

	if itemid != item.GetId() {
		t.Fatalf("Expected '%s', but got '%s'", itemid, item.GetId())
	}
}
