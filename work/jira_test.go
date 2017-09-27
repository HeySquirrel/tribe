package work

import (
	"testing"
)

func setupJira(t *testing.T) ItemServer {
	api, err := NewJiraFromConfig("rsjira")
	if err != nil {
		t.Fatal(err)
	}
	return api
}

func TestJiraGetItem(t *testing.T) {
	api := setupJira(t)
	expectedType := "Bug"

	item, err := api.GetItem("HIL-78")
	if err != nil {
		t.Fatal(err)
	}

	if expectedType != item.GetType() {
		t.Fatalf("Expected workitem to be of type '%s', but was: '%s'", expectedType, item.GetType())
	}
}

func TestJiraNotFoundItem(t *testing.T) {
	api := setupJira(t)
	itemid := "NOTID"

	item, err := api.GetItem(itemid)
	if err == nil || err != ItemNotFoundError(itemid) {
		t.Fatal("Expected ItemNotFoundError")
	}

	if itemid != item.GetId() {
		t.Fatalf("Expected '%s', but got '%s'", itemid, item.GetId())
	}
}

func TestJiraGetItemWhenNoLoginRequired(t *testing.T) {
	api, err := NewJiraFromConfig("jboss")
	if err != nil {
		t.Fatal(err)
	}
	expectedType := "Feature Request"

	item, err := api.GetItem("FURNACE-37")
	if err != nil {
		t.Fatal(err)
	}

	if expectedType != item.GetType() {
		t.Fatalf("Expected workitem to be of type '%s', but was: '%s'", expectedType, item.GetType())
	}
}
