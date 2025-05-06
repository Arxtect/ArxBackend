package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/cmd/api"
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/cmd/config"
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/cmd/migration_cmd"
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/cmd/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "ArxBackend",
	Short:        "ArxBackend",
	SilenceUsage: true,
	Long:         `ArxBackend`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New("parameter error")
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func tip() {
	usageStr := ` 🚀 Can use` + `-h` + ` View command`
	fmt.Printf("%s\n", usageStr)

}

func init() {
	rootCmd.AddCommand(api.StartCmd)
	rootCmd.AddCommand(version.StartCmd)
	rootCmd.AddCommand(config.StartCmd)
	rootCmd.AddCommand(migration_cmd.MigrateCmd)
}

// Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
