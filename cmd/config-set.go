package cmd

import (
    //"fmt"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/vrypan/fargo/config"
)

var configsetCmd = &cobra.Command{
    Use:   "set [parameter] [value]",
    Short: "Set a parameter",
    Long: `Examples:
fargo config set hub.host 192.168.1.1
fargo config set hub.port 2283
fargo config set hub.ssl false`,
    Run: config_set,
}

func config_set(cmd *cobra.Command, args []string) {
    config.Load()
    if len(args) != 2 {
        panic("Wrong number of arguments")
    }
    viper.Set(args[0], args[1])
    viper.WriteConfig()
}
func init() {
    configCmd.AddCommand(configsetCmd)
}