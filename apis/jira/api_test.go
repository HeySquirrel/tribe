package jira

import (
	"fmt"
	"testing"
)

func TestGetByFormattedIds(t *testing.T) {
	api := NewFromConfig("rsjira")

	artifacts, _ := api.GetWorkItem("HIL-78")

	fmt.Printf("HERE: %v\n", artifacts)
}
