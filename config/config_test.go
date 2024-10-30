package config

import (
	"testing"
	"github.com/spf13/viper"
)

func Test_Config(t *testing.T) {
	Load()
	// Access configuration values
	hub_host := viper.GetString("hub.host")
	hub_port := viper.GetString("hub.port")
	hub_ssl := viper.GetString("hub.ssl")
	display_cols := viper.GetString("display.cols")

	// Print configuration values
	t.Logf("Hub.host: %#v\n", hub_host)
	t.Logf("Hub.port: %#v\n", hub_port)
	t.Logf("Hub.ssl: %#v\n", hub_ssl)
	t.Logf("Display.cols: %#v\n", display_cols)
	
}
