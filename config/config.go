package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var FARGO_VERSION string

// Initialize configuration using Viper
func Load() string { // Load config and return config file path
	configDir, err := ConfigDir()
	if err != nil {
		log.Fatalf("Error getting config file: %v", err)
	}
	viper.SetEnvPrefix("FARGO")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName("config")
	viper.AddConfigPath(configDir)
	viper.SetConfigType("yaml")

	defaultDownload := "~/Downloads"
	if home, err := os.UserHomeDir(); err == nil {
		defaultDownload = filepath.Join(home, "Downloads")
	}

	defaults := map[string]interface{}{
		"hub.host":     "hoyt.farcaster.xyz",
		"hub.port":     "2283",
		"hub.ssl":      "true",
		"download.dir": defaultDownload,
		"get.count":    20,
		"cast.fid":     0,
		"cast.privkey": "",
		"cast.pubkey":  "",
		"db.ttlhours":  24,
	}
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Creating %s", filepath.Join(configDir, "config.yaml"))
			viper.SafeWriteConfig()
		} else {
			log.Fatalf("Error reading config file: %v", err)
		}
	}
	return viper.ConfigFileUsed()
}

var (
	GetString = viper.GetString
	GetInt    = viper.GetInt
	GetBool   = viper.GetBool
	BindPFlag = viper.BindPFlag
)
