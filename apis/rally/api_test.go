package rally

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"testing"
)

func TestGetByFormattedId(t *testing.T) {
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

	artifacts, _ := api.GetByFormattedId("S144101")

	fmt.Printf("HERE: %v\n", artifacts)
}
