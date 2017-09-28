// +build authenticated

package authenticatedtests

import "testing"
import "github.com/HeySquirrel/tribe/work"

func TestJiraGetItemWhenLoginRequired(t *testing.T) {
	api := work.SetupServer(t, "rsjira")
	expectedType := "Bug"

	item, err := api.GetItem("HIL-78")
	if err != nil {
		t.Fatal(err)
	}

	if expectedType != item.GetType() {
		t.Fatalf("Expected workitem to be of type '%s', but was: '%s'", expectedType, item.GetType())
	}
}
