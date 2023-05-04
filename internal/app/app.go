package app

import (
	"github.com/Spartan0nix/zabbix-map-builder-go/internal/api"
	zbxMap "github.com/Spartan0nix/zabbix-map-builder-go/internal/map"
)

// RunApp is used to run the main logic of the application.
func RunApp(file string, options *Options) error {
	// Retrieve the list of hosts mappings for the input file
	mappings, err := ReadInput(file)
	if err != nil {
		return err
	}

	// Initialize an api client.
	client, err := api.InitApi(options.ZabbixUrl, options.ZabbixUser, options.ZabbixPwd)
	if err != nil {
		return err
	}

	// Catch logout error
	defer func() {
		err = client.Logout()
	}()

	// Remove duplicate for the hosts mappings
	// Make it easier to retrieve id of each hosts
	hosts := getUniqueHosts(mappings)

	// Retrieve the id of each hosts and provide a mapping 'host' -> 'Zabbix id'
	hosts, err = api.GetHostsId(client, hosts)
	if err != nil {
		return err
	}

	// Construct map options
	mapOptions := zbxMap.MapOptions{
		Name:         options.Name,
		Color:        options.Color,
		TriggerColor: options.TriggerColor,
	}

	// Validate the options
	err = mapOptions.Validate()
	if err != nil {
		return err
	}

	// Build the map create request
	m, err := zbxMap.BuildMap(client, mappings, hosts, &mapOptions)
	if err != nil {
		return err
	}

	// Create the map using the previously build request
	err = zbxMap.CreateMap(client, m, options.OutFile)
	if err != nil {
		return err
	}

	// Allow to return errors from the defer function (API logout)
	return err
}
