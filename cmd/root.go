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
var StackHosts bool
var GlobalLogger *logging.Logger
var Debug bool

func init() {
	rootCmd.PersistentFlags().StringVar(&Name, "name", "", "name of map")
	rootCmd.PersistentFlags().StringVarP(&File, "file", "f", "", "input file")
	rootCmd.PersistentFlags().StringVarP(&OutFile, "output", "o", "", "output the parameters used to create the map")
	rootCmd.PersistentFlags().StringVarP(&Color, "color", "c", "000000", "color in hexadecimal used for the links between each hosts")
	rootCmd.PersistentFlags().StringVar(&TriggerColor, "trigger-color", "DD0000", "color in hexadecimal used for the links between each hosts when a trigger is in problem state")
	rootCmd.PersistentFlags().BoolVar(&StackHosts, "stack-hosts", false, "connect multiple links to a single host. If set to false, if mapping will have is own hosts. This can be useful for infrastructure with redundant connexion")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "v", false, "enable debug logging verbosity")

	rootCmd.MarkPersistentFlagRequired("name")
	rootCmd.MarkPersistentFlagRequired("file")

	// Init a new global logger
	GlobalLogger = logging.NewLogger(logging.Warning)
}

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Build a zabbix map using the given host mapping.",
	Long:  "This CLI tool is used to help administrator build a zabbix map using the given host mapping (network devices, etc.).",
	PreRun: func(cmd *cobra.Command, args []string) {
		// Check if the file flag was set correctly.
		checkRequiredFlag(Name, File)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Enable debug logger level.
		if Debug {
			GlobalLogger.Level = logging.Debug
		}

		// Retrieve the required environment variables.
		options, err := app.GetEnvironmentVariables()
		if err != nil {
			GlobalLogger.Error("error when reading the required environment variables", fmt.Sprintf("reason : %s", err))
			os.Exit(1)
		}

		options.Name = Name
		options.OutFile = OutFile
		options.Color = Color
		options.StackHosts = StackHosts
		options.TriggerColor = TriggerColor

		// Run the application.
		err = app.RunApp(File, options)
		if err != nil {
			GlobalLogger.Error("error when executing command", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Error when during command initialization.\nReason : %v", err)
	}
}

// checkRequiredFlag is used to validate the required flags.
func checkRequiredFlag(name string, file string) {
	if name == "" {
		GlobalLogger.Error("'name' flag is required and cannot be empty")
		os.Exit(1)
	}

	if file == "" {
		GlobalLogger.Error("'file' flag is required and cannot be empty")
		os.Exit(1)
	}

	if _, err := os.Stat(File); err != nil {
		GlobalLogger.Error(fmt.Sprintf("error while reading file '%s'.", File), err)
		os.Exit(1)
	}
}
