package config

import (
	"fmt"
	"log"
	"github.com/spf13/viper"
	"strings"
)

// Initialize configuration using Viper
func Load() {
	//viper.AutomaticEnv()
	viper.SetEnvPrefix("FARGO")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.fargo")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	
}

var GetString = viper.GetString
var GetBool = viper.GetBool

func Show() {
	Load()

	// Access configuration values
	appName := viper.GetString("app.name")
	appVersion := viper.GetString("app.version")
	serverHost := viper.GetString("server.host")
	serverPort := viper.GetInt("server.port")
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.password")
	dbName := viper.GetString("database.dbname")

	// Print configuration values
	fmt.Printf("App Name: %s\n", appName)
	fmt.Printf("App Version: %s\n", appVersion)
	fmt.Printf("Server Host: %s\n", serverHost)
	fmt.Printf("Server Port: %d\n", serverPort)
	fmt.Printf("Database User: %s\n", dbUser)
	fmt.Printf("Database Password: %s\n", dbPassword)
	fmt.Printf("Database Name: %s\n", dbName)
}
