package config

import (
	"fmt"
	"github.com/spf13/viper"
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

func RallyApiKey() string {
	return viper.GetString("rally.key")
}
