package cmd

import (
	"fmt"
	"os"

	"github.com/Spartan0nix/zabbix-map-builder-go/internal/app"
	"github.com/Spartan0nix/zabbix-map-builder-go/internal/logging"
	"github.com/spf13/cobra"
)

var Name string
var File string
var Color string
var TriggerColor string
var Width int
var Height int
var Spacer int
var StackHosts []bool
var DryRun bool

// checkCreateRequiredFlag is used to validate the required flags.
func checkCreateRequiredFlag(name string, file string) string {
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

// newCreateCmd is used to generate the create command for the CLI
func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a zabbix map on the server using the given host mapping.",
		Long:  "Create a zabbix map on the server using the given host mapping.",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Check if the file flag was set correctly.
			if err := checkCreateRequiredFlag(Name, File); err != "" {
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
			options.TriggerColor = TriggerColor
			options.Height = Height
			options.Width = Width
			options.Spacer = Spacer
			options.StackHosts = StackHosts[0]
			options.DryRun = DryRun

			// Run the application.
			err = app.RunCreate(File, options, GlobalLogger)
			if err != nil {
				GlobalLogger.Error("error when executing the command", err)
				os.Exit(1)
			}
		},
	}

	// Set all the persistent flag
	cmd.PersistentFlags().StringVar(&Name, "name", "", "name of the map")
	cmd.PersistentFlags().StringVarP(&File, "file", "f", "", "file containing the hosts mapping")
	cmd.PersistentFlags().StringVarP(&Color, "color", "c", "000000", "color in hexadecimal used for the links between each hosts")
	cmd.PersistentFlags().StringVar(&TriggerColor, "trigger-color", "DD0000", "color in hexadecimal used for the links between each hosts when a trigger is in problem state")
	cmd.PersistentFlags().IntVar(&Height, "height", 800, "height in pixel of the map")
	cmd.PersistentFlags().IntVar(&Width, "width", 800, "width in pixel of the map")
	cmd.PersistentFlags().IntVar(&Spacer, "spacer", 100, "space in pixel between each host (example : X_host2 = X_host1 + <value>)")
	cmd.PersistentFlags().BoolSliceVar(&StackHosts, "stack-hosts", []bool{true}, "connect multiple links to a single host. If set to false, each mapping will have is own hosts (local and remote). This can be useful for infrastructure with redundant connexion")
	cmd.PersistentFlags().BoolVar(&DryRun, "dry-run", false, "output to the shell the map definition without created it on the server")
	cmd.MarkPersistentFlagRequired("name")
	cmd.MarkPersistentFlagRequired("file")

	return cmd
}
