package rally

import (
	"fmt"
	"github.com/heysquirrel/tribe/config"
	"testing"
)

func TestGetByFormattedIds(t *testing.T) {
	apikey := config.RallyApiKey()
	api := New(apikey)

	artifacts, _ := api.GetByFormattedIds("S144101")

	fmt.Printf("HERE: %v\n", artifacts)
}
