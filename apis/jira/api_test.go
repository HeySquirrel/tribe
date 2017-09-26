package jira

import (
	"github.com/heysquirrel/tribe/apis"
	"testing"
)

func setup(t *testing.T) apis.WorkItemServer {
	api, err := NewFromConfig("rsjira")
	if err != nil {
		t.Fatal(err)
	}
	return api
}

func TestGetWorkItem(t *testing.T) {
	api := setup(t)
	expectedType := "Bug"

	item, err := api.GetWorkItem("HIL-78")
	if err != nil {
		t.Fatal(err)
	}

	if expectedType != item.GetType() {
		t.Fatalf("Expected workitem to be of type '%s', but was: '%s'", expectedType, item.GetType())
	}
}

func TestNotFoundWorkItem(t *testing.T) {
	api := setup(t)
	itemid := "NOTID"

	item, err := api.GetWorkItem(itemid)
	if err == nil && err != apis.ItemNotFoundError(itemid) {
		t.Fatal("Expected ItemNotFoundError")
	}

	if itemid != item.GetId() {
		t.Fatalf("Expected '%s', but got '%s'", itemid, item.GetId())
	}
}
