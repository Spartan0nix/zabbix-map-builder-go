package app

import (
	"encoding/json"
	"fmt"
	"os"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
	"github.com/Spartan0nix/zabbix-map-builder-go/internal/api"
	zbxmap "github.com/Spartan0nix/zabbix-map-builder-go/internal/map"
)

func outputToFile(file string, m *zabbixgosdk.MapCreateParameters) error {
	if file == "" {
		return fmt.Errorf("file name cannot be empty")
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// RunApp is used to run the main logic of the application.
func RunApp(file string, options *Options) error {
	// Retrieve the list of hosts mappings for the input file
	// fmt.Println("reading input file")
	mappings, err := ReadInput(file)
	if err != nil {
		return err
	}

	// Initialize an api client.
	// fmt.Println("initializing API client")
	client, err := api.InitApi(options.ZabbixUrl, options.ZabbixUser, options.ZabbixPwd)
	if err != nil {
		return err
	}

	// Catch logout error
	defer func() {
		err = client.Logout()
	}()

	// Remove duplicate from the hosts mappings and associate 'host' -> 'hostid'
	// Make it easier to retrieve id of each hosts
	// fmt.Println("retrieving the list of hosts")
	hosts, err := getUniqueHosts(client, mappings)
	if err != nil {
		return err
	}

	// Remove duplicate from the hosts mappings and associate 'image' -> 'imageid'
	// Make it easier to retrieve id of each hosts
	// fmt.Println("retrieving the list of images")
	images, err := getUniqueImages(client, mappings)
	if err != nil {
		return err
	}

	// Construct map options
	// fmt.Println("building map options")
	mapOptions := zbxmap.MapOptions{
		Name:         options.Name,
		Color:        options.Color,
		TriggerColor: options.TriggerColor,
		StackHosts:   options.StackHosts,
		Mappings:     mappings,
		Hosts:        hosts,
		Images:       images,
	}

	// Validate the options
	// fmt.Println("validating map options")
	err = mapOptions.Validate()
	if err != nil {
		return err
	}

	// Build the map create request
	// fmt.Println("build the map")
	m, err := zbxmap.BuildMap(client, &mapOptions)
	if err != nil {
		return err
	}

	// Store the create request if asked before executing it on the server
	if options.OutFile != "" {
		err = outputToFile(options.OutFile, m)
		if err != nil {
			return err
		}
	}

	// If dry-run was set to true, output the map definition to the shell
	if options.DryRun {
		// Convert the request parameters to a slice of byte before output the content as a string to the shell
		b, err := json.Marshal(m)
		if err != nil {
			return err
		}

		fmt.Println(string(b))

		return err
	}

	// Create the map using the previously build request
	err = zbxmap.CreateMap(client, m)
	if err != nil {
		return err
	}

	// Allow to return errors from the defer function (API logout)
	return err
}
