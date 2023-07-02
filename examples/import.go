package main

import (
	"fmt"
	"log"
	"os"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
	"github.com/spf13/cobra"
)

var importFile string

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Import a Zabbix export file.",
	Long:  "Import the given Zabbix export file.",
	Run: func(cmd *cobra.Command, args []string) {
		run(importFile)
	},
}

type env struct {
	Url  string
	User string
	Pwd  string
}

type importRequest struct {
	Format string                `json:"format"`
	Source string                `json:"source"`
	Rules  map[string]importRule `json:"rules"`
}

type importRule struct {
	CreateMissing  bool `json:"createMissing,omitempty"`
	UpdateExisting bool `json:"updateExisting,omitempty"`
	DeleteMissing  bool `json:"deleteMissing,omitempty"`
}

func init() {
	rootCmd.PersistentFlags().StringVar(&importFile, "file", "", "path of the Zabbix export file")
	rootCmd.MarkFlagRequired("file")
}

func getEnvironmentVariables() (*env, error) {
	env := &env{}

	env.Url = os.Getenv("ZABBIX_URL")
	if env.Url == "" {
		return nil, fmt.Errorf("missing required variables 'ZABBIX_URL'")
	}

	env.User = os.Getenv("ZABBIX_USER")
	if env.User == "" {
		return nil, fmt.Errorf("missing required variables 'ZABBIX_USER'")
	}

	env.Pwd = os.Getenv("ZABBIX_PWD")
	if env.Pwd == "" {
		return nil, fmt.Errorf("missing required variables 'ZABBIX_PWD'")
	}

	return env, nil
}

func run(file string) {
	env, err := getEnvironmentVariables()
	if err != nil {
		log.Fatal(err)
	}

	client := zabbixgosdk.NewZabbixService()
	client.SetUrl(env.Url)
	client.SetUser(&zabbixgosdk.ApiUser{
		User: env.User,
		Pwd:  env.Pwd,
	})

	b, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("error while reading file '%s'.\nReason : %v", file, err)
	}

	err = client.Authenticate()
	if err != nil {
		log.Fatalf("error while authenticating.\nReason : %v", err)
	}

	params := importRequest{
		Format: "json",
		Source: string(b),
		Rules: map[string]importRule{
			"groups": {
				CreateMissing:  true,
				UpdateExisting: true,
			},
			"hosts": {
				CreateMissing:  true,
				UpdateExisting: true,
			},
			"items": {
				CreateMissing:  true,
				UpdateExisting: true,
			},
			"triggers": {
				CreateMissing:  true,
				UpdateExisting: true,
			},
			"valueMaps": {
				CreateMissing:  true,
				UpdateExisting: true,
			},
		},
	}

	req := client.Host.Client.NewRequest("configuration.import", params)

	data, err := client.Host.Client.Post(req)
	if err != nil {
		log.Fatalf("error while executing import request.\nReason : %v", err)
	}

	var res bool
	err = client.HostGroup.Client.ConvertResponse(*data, &res)
	if err != nil {
		log.Fatalf("error while converting response to structure.\nReason : %v", err)
	}

	if !res {
		log.Fatalf("a non true result was returned.\nReturned : %v", res)
	}
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
