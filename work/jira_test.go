package work

import (
	"testing"
)

const UNAUTHENTICATED_JIRA = "jboss"

func TestJiraGetItem(t *testing.T) {
	api := SetupServer(t, UNAUTHENTICATED_JIRA)
	expectedType := "Feature Request"

	item, err := api.GetItem("FURNACE-37")
	if err != nil {
		t.Fatal(err)
	}

	if expectedType != item.GetType() {
		t.Fatalf("Expected workitem to be of type '%s', but was: '%s'", expectedType, item.GetType())
	}
}

func TestJiraNotFoundItem(t *testing.T) {
	api := SetupServer(t, UNAUTHENTICATED_JIRA)
	itemid := "NOTID"

	item, err := api.GetItem(itemid)
	if err == nil || err != ItemNotFoundError(itemid) {
		t.Fatal("Expected ItemNotFoundError")
	}

	if item != nil {
		t.Fatalf("Not found item should have been nil, got item with id '%s'", item.GetId())
	}
}
