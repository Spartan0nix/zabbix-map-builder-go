package cmd

import (
	"log"

	"github.com/Spartan0nix/zabbix-map-builder-go/internal/logging"
	"github.com/spf13/cobra"
)

var GlobalLogger *logging.Logger
var OutFile string
var Debug bool

func init() {
	// Init a new global logger
	GlobalLogger = logging.NewLogger(logging.Warning)
}

// newRootCmd is used to generate the root command for the CLI
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "zabbix-map-builder",
		Short:         "Build a zabbix map using the given host mapping.",
		Long:          "This CLI tool is used to help administrator build a zabbix map using the given host mappings (network devices, etc.).",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Set all the persistent flag
	cmd.PersistentFlags().BoolVarP(&Debug, "debug", "v", false, "enable debug logging verbosity")
	cmd.PersistentFlags().StringVarP(&OutFile, "output", "o", "", "output the parameters used to create the map to a file")
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newGenerateCmd())

	return cmd
}

func Execute() {
	rootCmd := newRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("error during command initialization.\nReason : %v", err)
	}
}
