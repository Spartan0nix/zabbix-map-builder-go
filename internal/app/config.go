package app

import (
	"encoding/json"
	"fmt"
	"os"

	zbxMap "github.com/Spartan0nix/zabbix-map-builder-go/internal/map"
)

type Options struct {
	ZabbixUrl    string
	ZabbixUser   string
	ZabbixPwd    string
	Name         string
	OutFile      string
	Color        string
	TriggerColor string
	Height       string
	Width        string
	StackHosts   bool
	DryRun       bool
}

// GetEnvironmentVariables is used to retrive the required environment variables for the Zabbix API.
func GetEnvironmentVariables() (*Options, error) {
	vars := Options{}

	if vars.ZabbixUrl = os.Getenv("ZABBIX_URL"); vars.ZabbixUrl == "" {
		return nil, fmt.Errorf("required environment variable 'ZABBIX_URL' is not set")
	}

	if vars.ZabbixUser = os.Getenv("ZABBIX_USER"); vars.ZabbixUser == "" {
		return nil, fmt.Errorf("required environment variable 'ZABBIX_USER' is not set")
	}

	if vars.ZabbixPwd = os.Getenv("ZABBIX_PWD"); vars.ZabbixPwd == "" {
		return nil, fmt.Errorf("required environment variable 'ZABBIX_PWD' is not set")
	}

	return &vars, nil
}

// readInput is used to read data from the given file and return a list of Host.
func ReadInput(file string) ([]*zbxMap.Mapping, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	entries := make([]*zbxMap.Mapping, 0)
	err = json.Unmarshal(b, &entries)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no mapping were found in '%s'", file)
	}

	return entries, nil
}
