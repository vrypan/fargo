package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/vrypan/fargo/config"
)

var configgetCmd = &cobra.Command{
    Use:   "get [parameter1] [parameter2] ...",
    Short: "Get parameter values",
    Long: `Examples:
fargo config get hub.host
fargo config get hub.host hub.port hub.ssl`,
    Run: config_get,
}

func config_get(cmd *cobra.Command, args []string) {
    config.Load()
    for _,arg := range args {
        fmt.Printf("%s: %s\n", arg, viper.GetString(arg) )
    }
}
func init() {
    configCmd.AddCommand(configgetCmd)
}