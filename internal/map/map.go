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
	StackHosts   bool
	Mappings     []*Mapping
	Hosts        map[string]string
	Images       map[string]string
}

// Validate is used to validate options that will be passed to a map.
func (o *MapOptions) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("a name is required to create the map, used the 'name' flag to set one")
	}

	if o.Color == "" {
		o.Color = "000000"
	}

	if o.TriggerColor == "" {
		o.TriggerColor = "DD0000"
	}

	if o.Color != "000000" {
		if err := validateHexa(o.Color); err != nil {
			return err
		}
	}

	if o.TriggerColor != "DD0000" {
		if err := validateHexa(o.TriggerColor); err != nil {
			return err
		}
	}

	if o.Mappings == nil {
		return fmt.Errorf("no mappings were passed to the build function")
	}

	if o.Hosts == nil {
		return fmt.Errorf("no mapping 'host' -> 'hostid' was passed to the build function")
	}

	if o.Images == nil {
		return fmt.Errorf("no images were passed to the build function")
	}

	return nil
}

// buildElementId is used to update the local and remote elements id based on the number of hosts that already exists with the same id.
// Used only is --stack-hosts is set to false.
func buildElementsId(zbxMap *zabbixgosdk.MapCreateParameters, localElementId string, remoteElementId string) (string, string) {
	localCount := 0
	remoteCount := 0

	for _, element := range zbxMap.Elements {
		elementId := element.Elements.([]zabbixgosdk.MapElementHost)[0].Id

		if elementId == localElementId {
			localCount++
		}
		if elementId == remoteElementId {
			remoteCount++
		}
	}

	if localCount > 0 {
		localElementId = fmt.Sprintf("%s-%d", localElementId, localCount+1)
	}

	if remoteCount > 0 {
		remoteElementId = fmt.Sprintf("%s-%d", remoteElementId, remoteCount+1)
	}

	return localElementId, remoteElementId
}

// BuildMap is used to build a map with the given mapping.
func BuildMap(client *zabbixgosdk.ZabbixService, options *MapOptions) (*zabbixgosdk.MapCreateParameters, error) {
	zbxMap := &zabbixgosdk.MapCreateParameters{}
	zbxMap.Name = options.Name
	zbxMap.Width = "800"
	zbxMap.Height = "800"

	// Loop over each mapping
	for _, mapping := range options.Mappings {
		localElementId := options.Hosts[mapping.LocalHost]
		remoteElementId := options.Hosts[mapping.RemoteHost]

		// If hosts should not be stacked, update the elementsId by appending '-<number-of-element-already-present + 1>'
		if !options.StackHosts {
			localElementId, remoteElementId = buildElementsId(zbxMap, localElementId, remoteElementId)
		}

		// Add the hosts to the map
		zbxMap = addHosts(zbxMap, localElementId, options.Hosts[mapping.LocalHost], options.Images[mapping.LocalImage])
		zbxMap = addHosts(zbxMap, remoteElementId, options.Hosts[mapping.RemoteHost], options.Images[mapping.RemoteImage])

		// Retriev the triggers id based on the given pattern for each hosts
		localTriggerId, err := getTriggerId(client, options.Hosts[mapping.LocalHost], mapping.LocalTriggerPattern)
		if err != nil {
			return nil, err
		}

		remoteTriggerId, err := getTriggerId(client, options.Hosts[mapping.RemoteHost], mapping.RemoteTriggerPattern)
		if err != nil {
			return nil, err
		}

		// Add the link to the map
		zbxMap = addLink(zbxMap, &linkParameters{
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
