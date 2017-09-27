package config

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
)

func init() {
	viper.SetConfigName(".tribe") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func SetConfigFile(cfgFile string) {
	viper.SetConfigFile(cfgFile)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

type ServerName string

func ItemServer(servername ServerName) map[string]string {
	key := fmt.Sprintf("workitemservers.%s", servername)
	return viper.GetStringMapString(key)
}

func ItemServers() []ServerName {
	servers := viper.GetStringMap("workitemservers")
	keys := reflect.ValueOf(servers).MapKeys()
	names := make([]ServerName, len(keys))

	for i := 0; i < len(keys); i++ {
		names[i] = ServerName(keys[i].String())
	}

	return names
}
