package config

import (
	"log"
	"github.com/spf13/viper"
	"strings"
)

var FARGO_VERSION string

// Initialize configuration using Viper
func Load() string { // Load config and return config file path
	//viper.AutomaticEnv()
	viper.SetEnvPrefix("FARGO")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.fargo")
	viper.SetConfigType("yaml")

	viper.SetDefault("hub.host", "hoyt.farcaster.xyz")
	viper.SetDefault("hub.port", "2283")
	viper.SetDefault("hub.ssl", "true")
	viper.SetDefault("downloads.dir", "~/Downloads")

	if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            log.Println("Creating ~/.fargo/config.yaml")
            viper.SafeWriteConfig()
        } else {
            log.Fatalf("Error reading config file: %v", err)
        }
    }
    return viper.ConfigFileUsed()
}

var GetString = viper.GetString
var GetBool = viper.GetBool
