package work

import (
	"testing"

	"github.com/HeySquirrel/tribe/config"
	"github.com/spf13/viper"
)

func SetupServer(t *testing.T, servername string) ItemServer {
	api, err := NewItemServerFromConfig(config.ServerName(servername))
	if err != nil {
		t.Fatalf("Requested Server: '%s', From Config: %s - %v", servername, viper.ConfigFileUsed(), err)
	}
	return api
}
