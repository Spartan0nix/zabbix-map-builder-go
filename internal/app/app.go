package app

import (
	"github.com/Spartan0nix/zabbix-map-builder-go/internal/api"
	zbxMap "github.com/Spartan0nix/zabbix-map-builder-go/internal/map"
)

// RunApp is used to run the main logic of the application.
func RunApp(file string, env *Env, outFile string) error {
	// Retrieve the list of hosts mappings for the input file
	mappings, err := ReadInput(file)
	if err != nil {
		return err
	}

	// Initialize an api client.
	client, err := api.InitApi(env.ZabbixUrl, env.ZabbixUser, env.ZabbixPwd)
	if err != nil {
		return err
	}

	// Remove duplicate for the hosts mappings
	// Make it easier to retrieve id of each hosts
	hosts := getUniqueHosts(mappings)

	// Retrieve the id of each hosts and provide a mapping 'host' -> 'Zabbix id'
	hosts, err = api.GetHostsId(client, hosts)
	if err != nil {
		return err
	}

	m, err := zbxMap.BuildMap(client, mappings, hosts)
	if err != nil {
		return err
	}

	err = zbxMap.CreateMap(client, m, outFile)
	if err != nil {
		return err
	}

	return nil
}
