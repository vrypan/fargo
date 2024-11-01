package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/vrypan/fargo/config"
)

var configlsCmd = &cobra.Command{
    Use:   "ls",
    Short: "Get a list of all parameters and their values",
    Run: config_ls,
}

func traverse(parent string, data map[string]interface{}) {
    for key, value := range data {
        prefix := parent + "."
        if parent == "" {
            prefix = ""
        }
        switch v := value.(type) {
        case string:
            fmt.Printf("%s: %s\n", prefix+key, v)
        case int:
            fmt.Printf("%s: %d\n", prefix+key, v)
        case []interface{}:
            fmt.Printf("%s:", prefix+key)
            for _, item := range v {
                fmt.Printf(" %v", item)
            }
            fmt.Printf("\n")
        case map[string]interface{}:
            traverse(prefix+key, v) // Recursive call to traverse nested map
        default:
            panic("Unknown value")
        }
    }
}
func config_ls(cmd *cobra.Command, args []string) {
        config.Load()
        settings := viper.AllSettings()
        traverse("",settings)
        
    }
func init() {
    configCmd.AddCommand(configlsCmd)
}