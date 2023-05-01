package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Spartan0nix/zabbix-map-builder-go/internal/app"
	"github.com/Spartan0nix/zabbix-map-builder-go/internal/logging"
	"github.com/spf13/cobra"
)

var File string
var OutFile string
var GlobalLogger *logging.Logger
var Debug bool

func init() {
	rootCmd.PersistentFlags().StringVarP(&File, "file", "f", "", "input file")
	rootCmd.PersistentFlags().StringVarP(&OutFile, "output", "o", "", "output the parameters used to create the map")
	// Init a new global logger
	GlobalLogger = logging.NewLogger(logging.Warning)
}

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Build a zabbix map using the given host mapping.",
	Long:  "This CLI tool is used to help administrator build a zabbix map using the given host mapping (network devices, etc.).",
	Run: func(cmd *cobra.Command, args []string) {
		// Enable debug logger level.
		if Debug {
			GlobalLogger.Level = logging.Debug
		}

		// Check if the file flag was set correctly.
		checkFileFlag(File)

		// Retrieve the required environment variables.
		env, err := app.GetEnvironmentVariables()
		if err != nil {
			GlobalLogger.Error("error when reading the required environment variables", fmt.Sprintf("reason : %s", err))
			os.Exit(1)
		}

		// Run the application.
		err = app.RunApp(File, env, OutFile)
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

// checkFileFlag is used to validate the given file variable.
func checkFileFlag(file string) {
	if file == "" {
		GlobalLogger.Error("flag '--file' cannot be empty")
		os.Exit(1)
	}

	if _, err := os.Stat(File); err != nil {
		GlobalLogger.Error(fmt.Sprintf("error while reading file '%s'.", File), err)
		os.Exit(1)
	}
}
