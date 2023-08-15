package _map

import (
	"fmt"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

// Mapping define the properties used to create an hosts mapping on a Zabbix map.
type Mapping struct {
	LocalHost            string `json:"local_host"`
	LocalInterface       string `json:"local_interface,omitempty"`
	LocalTriggerPattern  string `json:"local_trigger_pattern"`
	LocalImage           string `json:"local_image"`
	RemoteHost           string `json:"remote_host"`
	RemoteInterface      string `json:"remote_interface,omitempty"`
	RemoteTriggerPattern string `json:"remote_trigger_pattern"`
	RemoteImage          string `json:"remote_image"`
}

// MapOptions define the available options that can be passed to customize the map rendering.
type MapOptions struct {
	Name         string
	Color        string
	TriggerColor string
	Height       int
	Width        int
	Spacer       int
	StackHosts   bool
	Mappings     []*Mapping
	Hosts        map[string]string
	Images       map[string]string
}

// Validate is used to validate options that will be passed to a map.
func (o *MapOptions) Validate() error {
	var err error

	if o.Name == "" {
		return fmt.Errorf("a name is required to create the map, used the 'name' flag to set one")
	}

	if o.Color == "" {
		o.Color = "000000"
	} else {
		err = validateHexa(o.Color)
		if err != nil {
			return err
		}
	}

	if o.TriggerColor == "" {
		o.TriggerColor = "DD0000"
	} else {
		err = validateHexa(o.TriggerColor)
		if err != nil {
			return err
		}
	}

	if o.Mappings == nil {
		err = fmt.Errorf("no mappings were passed to the build function")
		if err != nil {
			return err
		}
	}

	if o.Hosts == nil {
		err = fmt.Errorf("no mapping 'host' -> 'hostid' was passed to the build function")
		if err != nil {
			return err
		}
	}

	if o.Images == nil {
		err = fmt.Errorf("no images were passed to the build function")
		if err != nil {
			return err
		}
	}

	return err
}

func BuildMap(client *zabbixgosdk.ZabbixService, options *MapOptions) (*zabbixgosdk.MapCreateParameters, error) {
	var unstackedHosts = make(map[string]int8, 0)

	zbxMap := &zabbixgosdk.MapCreateParameters{}
	zbxMap.Name = options.Name
	zbxMap.Height = fmt.Sprintf("%d", options.Height)
	zbxMap.Width = fmt.Sprintf("%d", options.Width)

	position, err := initPosition(options.Width, options.Height, options.Spacer)
	if err != nil {
		return nil, err
	}

	// Loop over each mapping
	for _, mapping := range options.Mappings {
		localElementId := options.Hosts[mapping.LocalHost]
		remoteElementId := options.Hosts[mapping.RemoteHost]

		// If hosts should not be stacked, update the elementsId by appending '-<number-of-element-already-present + 1>'
		if !options.StackHosts {
			// Increment the number of hosts
			unstackedHosts[mapping.LocalHost] += 1
			unstackedHosts[mapping.RemoteHost] += 1
			// Build the new elements ids
			localElementId = fmt.Sprintf("%s-%d", options.Hosts[mapping.LocalHost], unstackedHosts[mapping.LocalHost])
			remoteElementId = fmt.Sprintf("%s-%d", options.Hosts[mapping.RemoteHost], unstackedHosts[mapping.LocalHost])
		}

		// Add the local host to the map
		zbxMap = addHosts(zbxMap, &hostParameters{
			id:       localElementId,
			name:     options.Hosts[mapping.LocalHost],
			image:    options.Images[mapping.LocalImage],
			position: position,
		})

		// Add the remote host to the map
		zbxMap = addHosts(zbxMap, &hostParameters{
			id:       remoteElementId,
			name:     options.Hosts[mapping.RemoteHost],
			image:    options.Images[mapping.RemoteImage],
			position: position,
		})

		// Retrieve the triggers id based on the given pattern for each hosts
		localTriggerId, err := getTriggerId(client, options.Hosts[mapping.LocalHost], mapping.LocalTriggerPattern)
		if err != nil {
			return nil, err
		}

		remoteTriggerId, err := getTriggerId(client, options.Hosts[mapping.RemoteHost], mapping.RemoteTriggerPattern)
		if err != nil {
			return nil, err
		}

		// Add the link to the map
		addLink(zbxMap, linkParameters{
			localElement:     localElementId,
			localTrigger:     localTriggerId,
			remoteElement:    remoteElementId,
			remoteTrigger:    remoteTriggerId,
			linkColor:        options.Color,
			triggerLinkColor: options.TriggerColor,
		})
	}

	return zbxMap, nil
}

// CreateMap is used to create the given map.
// The map create parameters can also be exported to a file if a file path is specified.
func CreateMap(client *zabbixgosdk.ZabbixService, m *zabbixgosdk.MapCreateParameters) error {
	res, err := client.Map.Create(m)
	if err != nil {
		return err
	}

	if res == nil || len(res.MapIds) == 0 {
		return fmt.Errorf("an empty response was returned when creating the map")
	}

	return nil
}
