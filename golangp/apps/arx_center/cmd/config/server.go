package config

import (
	"fmt"
	"github.com/Arxtect/ArxBackend/golangp/config"
	"github.com/spf13/cobra"
	"log"
)

var (
	configYml string
	StartCmd  = &cobra.Command{
		Use:     "config",
		Short:   "Get Application config info",
		Example: "ArxBackend config -c config/settings-dev.yml",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings-dev.yml", "Start server with provided configuration file")
}

func run() {
	log.Println("🚗 Load configuration file ...")
	err := config.LoadEnv(configYml)
	if err != nil {
		log.Println("🚀 Load failed", err)
		return
	}
	fmt.Println("🚗 Load success.....", config.Env.Mode)
}
