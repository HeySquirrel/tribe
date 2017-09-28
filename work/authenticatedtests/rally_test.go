// +build authenticated

package authenticatedtests

import (
	"testing"

	"github.com/HeySquirrel/tribe/work"
)

func TestRallyGetItem(t *testing.T) {
	api := work.SetupServer(t, "rally1")
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
	api := work.SetupServer(t, "rally1")
	itemid := "NOTID"

	item, err := api.GetItem(itemid)
	if err == nil || err != work.ItemNotFoundError(itemid) {
		t.Fatal("Expected ItemNotFoundError")
	}

	if item != nil {
		t.Fatalf("Not found item should have been nil, got item with id '%s'", item.GetId())
	}
}
