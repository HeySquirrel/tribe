package work

import (
	"testing"
)

func setupRally(t *testing.T) ItemServer {
	api, err := NewRallyFromConfig("rally1")
	if err != nil {
		t.Fatal(err)
	}
	return api
}

func TestRallyGetItem(t *testing.T) {
	api := setupRally(t)
	expectedType := "HierarchicalRequirement"

	item, err := api.GetItem("S113541")
	if err != nil {
		t.Fatal(err)
	}

	if expectedType != item.GetType() {
		t.Fatalf("Expected workitem to be of type '%s', but was: '%s'", expectedType, item.GetType())
	}
}

func TestRallyNotFoundItem(t *testing.T) {
	api := setupRally(t)
	itemid := "NOTID"

	item, err := api.GetItem(itemid)
	if err == nil || err != ItemNotFoundError(itemid) {
		t.Fatal("Expected ItemNotFoundError")
	}

	if item != nil {
		t.Fatalf("Not found item should have been nil, got item with id '%s'", item.GetId())
	}
}
