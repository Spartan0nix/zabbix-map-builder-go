package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Spartan0nix/zabbix-map-builder-go/internal/app"
	"github.com/Spartan0nix/zabbix-map-builder-go/internal/logging"
	"github.com/spf13/cobra"
)

var Name string
var File string
var OutFile string
var Color string
var TriggerColor string
var StackHosts []bool
var GlobalLogger *logging.Logger
var Debug bool
var DryRun bool

func init() {
	// Init a new global logger
	GlobalLogger = logging.NewLogger(logging.Warning)
}

// newRootCmd is used to generate the root command for the CLI
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "",
		Short: "Build a zabbix map using the given host mapping.",
		Long:  "This CLI tool is used to help administrator build a zabbix map using the given host mappings (network devices, etc.).",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Check if the file flag was set correctly.
			if err := checkRequiredFlag(Name, File); err != "" {
				GlobalLogger.Error(err)
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Enable debug logger level.
			if Debug {
				GlobalLogger.Level = logging.Debug
			}

			// Retrieve the required environment variables.
			GlobalLogger.Debug("retrieving environment variables")
			options, err := app.GetEnvironmentVariables()
			if err != nil {
				GlobalLogger.Error("error when reading the required environment variables", fmt.Sprintf("reason : %s", err))
				os.Exit(1)
			}

			GlobalLogger.Debug(fmt.Sprintf("using the following environment variables :\nZABBIX_URL => %s\nZABBIX_USER => %s\nZABBIX_PWD => <masked-for-security-reason>", options.ZabbixUrl, options.ZabbixUser))

			options.Name = Name
			options.OutFile = OutFile
			options.Color = Color
			options.StackHosts = StackHosts[0]
			options.TriggerColor = TriggerColor
			options.DryRun = DryRun

			// Run the application.
			err = app.RunApp(File, options, GlobalLogger)
			if err != nil {
				GlobalLogger.Error("error when executing the command", err)
				os.Exit(1)
			}
		},
	}

	// Set all the persistent flag
	cmd.PersistentFlags().StringVar(&Name, "name", "", "name of map")
	cmd.PersistentFlags().StringVarP(&File, "file", "f", "", "input file")
	cmd.PersistentFlags().StringVarP(&OutFile, "output", "o", "", "output the parameters used to create the map")
	cmd.PersistentFlags().StringVarP(&Color, "color", "c", "000000", "color in hexadecimal used for the links between each hosts")
	cmd.PersistentFlags().StringVar(&TriggerColor, "trigger-color", "DD0000", "color in hexadecimal used for the links between each hosts when a trigger is in problem state")
	cmd.PersistentFlags().BoolSliceVar(&StackHosts, "stack-hosts", []bool{true}, "connect multiple links to a single host. If set to false, if mapping will have is own hosts. This can be useful for infrastructure with redundant connexion")
	cmd.PersistentFlags().BoolVarP(&Debug, "debug", "v", false, "enable debug logging verbosity")
	cmd.PersistentFlags().BoolVar(&DryRun, "dry-run", false, "output to the shell the map definition without created it on the server")
	cmd.MarkPersistentFlagRequired("name")
	cmd.MarkPersistentFlagRequired("file")

	return cmd
}

func Execute() {
	rootCmd := newRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("error during command initialization.\nReason : %v", err)
	}
}

// checkRequiredFlag is used to validate the required flags.
func checkRequiredFlag(name string, file string) string {
	if name == "" {
		return "'name' flag is required and cannot be empty"
	}

	if file == "" {
		return "'file' flag is required and cannot be empty"
	}

	if _, err := os.Stat(file); err != nil {
		return fmt.Sprintf("error while reading file '%s'.", file)
	}

	return ""
}
