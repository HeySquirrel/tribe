package rally

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"testing"
)

func TestGetByFormattedIds(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	configFile := filepath.Join(usr.HomeDir, ".tribe")
	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		t.Fatal(err)
	}

	api := New(string(config))

	artifacts, _ := api.GetByFormattedIds("S144101")

	fmt.Printf("HERE: %v\n", artifacts)
}
