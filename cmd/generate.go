package cmd

import (
	"fmt"
	"os"

	"github.com/Spartan0nix/zabbix-map-builder-go/internal/app"
	"github.com/Spartan0nix/zabbix-map-builder-go/internal/logging"
	"github.com/spf13/cobra"
)

var Host string
var Port uint16
var Community string
var TriggerPattern string
var LocalImage string
var RemoteImage string

// checkGenerateRequiredFlag is used to validate the required flags.
func checkGenerateRequiredFlag(host string, community string) string {
	if host == "" {
		return "'host' flag is required and cannot be empty"
	}

	if community == "" {
		return "'community' flag is required and cannot be empty"
	}

	return ""
}

// newGenerateCmd is used to generate the generate command for the CLI
func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a mapping file for a given host.",
		Long:  "Generate a mapping file for a given host using SNMP and the host CDP cache.",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Check if the file flag was set correctly.
			if err := checkGenerateRequiredFlag(Host, Community); err != "" {
				GlobalLogger.Error(err)
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Enable debug logger level.
			if Debug {
				GlobalLogger.Level = logging.Debug
			}

			GlobalLogger.Debug(fmt.Sprintf("using the following environment variables :\nHost => %s\nCommunity => %s", Host, Community))

			options := app.GenerateOptions{
				Host:           Host,
				Community:      Community,
				Port:           Port,
				OutFile:        OutFile,
				TriggerPattern: TriggerPattern,
				LocalImage:     LocalImage,
				RemoteImage:    RemoteImage,
			}

			// Run the application.
			err := app.RunGenerate(&options, GlobalLogger)
			if err != nil {
				GlobalLogger.Error("error when executing the command", err)
				os.Exit(1)
			}
		},
	}

	triggerPatternDescription := `pattern used to populate fields 'local_trigger_pattern' and 'remote_trigger_pattern'.
Variable '#INTERFACE' can be utilize in the expression to reference the discovered interface.
Example :
- Pattern : "Operation Status of #INTERFACE"
- Data discovered using the SNMP CdpCacheTable : router1 [local] / eth0 -> router2 [remote] / eth1
- Resultant patterns :
	- local_trigger_pattern  : "Operation Status of eth0"
	- remote_trigger_pattern : "Operation Status of eth1"`

	cmd.PersistentFlags().StringVar(&Host, "host", "", "IP or DNS of the host")
	cmd.PersistentFlags().StringVarP(&Community, "community", "c", "", "SNMP community used to retrieve data")
	cmd.PersistentFlags().Uint16VarP(&Port, "port", "p", 161, "port to use for the SNMP requests.")
	cmd.PersistentFlags().StringVar(&TriggerPattern, "trigger-pattern", "", triggerPatternDescription)
	cmd.PersistentFlags().StringVar(&LocalImage, "local-host-image", "Switch_(64)", "name of the image used to populate field 'local_image'")
	cmd.PersistentFlags().StringVar(&RemoteImage, "remote-host-image", "Switch_(64)", "name of the image used to populate field 'remote_image")
	cmd.MarkPersistentFlagRequired("host")
	cmd.MarkPersistentFlagRequired("community")

	return cmd
}
