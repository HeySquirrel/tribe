package rally

import (
	"fmt"
	"testing"
)

func TestGetByFormattedIds(t *testing.T) {
	api, err := NewFromConfig("rally1")
	if err != nil {
		t.Fatal(err)
	}

	artifacts, _ := api.GetWorkItem("S113541")

	fmt.Printf("HERE: %v\n", artifacts)
}
