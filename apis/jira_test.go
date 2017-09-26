package apis

import (
	"testing"
)

func setupJira(t *testing.T) WorkItemServer {
	api, err := NewJiraFromConfig("rsjira")
	if err != nil {
		t.Fatal(err)
	}
	return api
}

func TestJiraGetWorkItem(t *testing.T) {
	api := setupJira(t)
	expectedType := "Bug"

	item, err := api.GetWorkItem("HIL-78")
	if err != nil {
		t.Fatal(err)
	}

	if expectedType != item.GetType() {
		t.Fatalf("Expected workitem to be of type '%s', but was: '%s'", expectedType, item.GetType())
	}
}

func TestJiraNotFoundWorkItem(t *testing.T) {
	api := setupJira(t)
	itemid := "NOTID"

	item, err := api.GetWorkItem(itemid)
	if err == nil || err != ItemNotFoundError(itemid) {
		t.Fatal("Expected ItemNotFoundError")
	}

	if itemid != item.GetId() {
		t.Fatalf("Expected '%s', but got '%s'", itemid, item.GetId())
	}
}
