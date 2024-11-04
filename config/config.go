package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
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
	viper.SetDefault("download.dir", "~/Downloads")
	viper.SetDefault("get.count", 20)
	viper.SetDefault("cast.fid", 0)
	viper.SetDefault("cast.privkey", "")
	viper.SetDefault("cast.pubkey", "")

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
var GetInt = viper.GetInt
var GetBool = viper.GetBool
var BindPFlag = viper.BindPFlag
