package app

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Spartan0nix/zabbix-map-builder-go/internal/api"
	"github.com/Spartan0nix/zabbix-map-builder-go/internal/logging"
	zbxmap "github.com/Spartan0nix/zabbix-map-builder-go/internal/map"
	"github.com/Spartan0nix/zabbix-map-builder-go/internal/snmp"
	"github.com/gosnmp/gosnmp"
)

// RunCreate is used to run the main logic for the create command.
func RunCreate(file string, options *MapOptions, logger *logging.Logger) error {
	if logger == nil {
		logger = logging.NewLogger(logging.Warning)
	}

	// Retrieve the list of hosts mappings for the input file
	logger.Debug(fmt.Sprintf("reading input file '%s'", file))
	mappings, err := ReadInput(file)
	if err != nil {
		return err
	}

	// Initialize an api client.
	logger.Debug("initializing the API client")
	client, err := api.InitApi(options.ZabbixUrl, options.ZabbixUser, options.ZabbixPwd)
	if err != nil {
		return err
	}

	// Catch logout error
	defer func() {
		client.Logout()
	}()

	// Remove duplicate from the hosts mappings and associate 'host' -> 'hostid'
	// Make it easier to retrieve id of each hosts
	logger.Debug("retrieving hosts information from the server")
	hosts, err := getUniqueHosts(client, mappings)
	if err != nil {
		return err
	}

	// Remove duplicate from the hosts mappings and associate 'image' -> 'imageid'
	// Make it easier to retrieve id of each hosts
	logger.Debug("retrieving images information from the server")
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
		Height:       options.Height,
		Width:        options.Width,
		Spacer:       options.Spacer,
		StackHosts:   options.StackHosts,
		Mappings:     mappings,
		Hosts:        hosts,
		Images:       images,
	}

	// Validate the options
	logger.Debug("validating the map configuration options")
	err = mapOptions.Validate()
	if err != nil {
		return err
	}

	if logger.Level >= logging.Debug {
		mapInfo := fmt.Sprintf("Name : %s\nLink color : %s\nTrigger color : %s\nStacked hosts : %t", mapOptions.Name, mapOptions.Color, mapOptions.TriggerColor, mapOptions.StackHosts)
		logger.Debug(fmt.Sprintf("the following options will be used to build the map :\n%s", mapInfo))
	}

	// Build the map create request
	logger.Debug("building the map")
	m, err := zbxmap.BuildMap(client, &mapOptions)
	if err != nil {
		return err
	}

	// Store the create request if asked before executing it on the server
	if options.OutFile != "" {
		logger.Debug(fmt.Sprintf("outputting the create request to '%s'", options.OutFile))

		b, err := json.Marshal(m)
		if err != nil {
			return err
		}

		err = outputToFile(options.OutFile, b)
		if err != nil {
			return err
		}
	} else {
		logger.Debug("'--ouptput' flag not used, skipping step.")
	}

	// If dry-run was set to true, output the map definition to the shell
	if options.DryRun {
		// Convert the request parameters to a slice of byte before output the content as a string to the shell
		logger.Debug("outputting map to the shell")
		b, err := json.Marshal(m)
		if err != nil {
			return err
		}

		fmt.Println(string(b))

		return err
	} else {
		logger.Debug("'--dry-run' flag not used, skipping step.")
	}

	// Create the map using the previously build request
	logger.Debug("creating the map on the server")
	err = zbxmap.CreateMap(client, m)
	if err != nil {
		return err
	}

	// Allow to return errors from the defer function (API logout)
	logger.Debug("all steps have been passed already, starting the exit process.")
	return err
}

// RunGenerate is used to run the main logic for the generate command.
func RunGenerate(options *GenerateOptions, logger *logging.Logger) error {
	if logger == nil {
		logger = logging.NewLogger(logging.Warning)
	}

	params := &gosnmp.GoSNMP{
		Target:    options.Host,
		Port:      options.Port,
		Community: options.Community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
	}

	err := params.Connect()
	if err != nil {
		return err
	}
	defer params.Conn.Close()

	// Retrieve the full cdpCacheTable
	oid := "1.3.6.1.4.1.9.9.23.1.2.1.1"
	res, err := snmp.WalkBulk(params, oid)
	if err != nil {
		return err
	}

	snmp.LogRequestDuration(logger, res.Duration)

	// Parse the cdpCacheTable to friendlier format
	cdpEntries := snmp.ParseCdpCache(res.Entries, logger)
	if len(cdpEntries) == 0 {
		return fmt.Errorf("no cdp data found on host '%s', check if cdp is up and running on the host", options.Host)
	}

	// Retrieve each local interface name
	err = snmp.GetLocalInterfacesName(params, cdpEntries, logger)
	if err != nil {
		return err
	}

	// Retrieve the local hostname
	localHostname, err := snmp.GetHostname(params, logger)
	if err != nil {
		return err
	}

	// Generate mappings
	mapping := generateMapping(cdpEntries, localHostname, &mappingOptions{
		TriggerPattern: options.TriggerPattern,
		LocalImage:     options.LocalImage,
		RemoteImage:    options.RemoteImage,
	})

	// Marshall data
	b, err := json.Marshal(&mapping)
	if err != nil {
		return err
	}

	if options.OutFile != "" {
		// Output to a file
		err = outputToFile(options.OutFile, b)
		if err != nil {
			return nil
		}
	} else {
		// Output to the shell
		fmt.Println(string(b))
	}

	return nil
}
