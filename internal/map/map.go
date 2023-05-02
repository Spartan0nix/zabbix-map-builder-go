package _map

import (
	"encoding/json"
	"fmt"
	"os"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

// Mapping define the properties used to create an hosts mapping on a Zabbix map.
type Mapping struct {
	LocalHost            string `json:"local_host"`
	LocalInterface       string `json:"local_interface"`
	LocalTriggerPattern  string `json:"local_trigger_pattern"`
	RemoteHost           string `json:"remote_host"`
	RemoteInterface      string `json:"remote_interface"`
	RemoteTriggerPattern string `json:"remote_trigger_pattern"`
}

// MapOptions define the available options that can be passed to customize the map rendering.
type MapOptions struct {
	Name         string
	Color        string
	TriggerColor string
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
		o.Color = "DD0000"
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

	return nil
}

// BuildMap is used to build a map with the given mapping.
func BuildMap(client *zabbixgosdk.ZabbixService, mappings []*Mapping, hosts map[string]string, options *MapOptions) (*zabbixgosdk.MapCreateParameters, error) {
	zbxMap := &zabbixgosdk.MapCreateParameters{}
	zbxMap.Name = options.Name
	zbxMap.Width = "800"
	zbxMap.Height = "800"
	var err error

	for _, mapping := range mappings {
		zbxMap = addHosts(zbxMap, mapping, hosts)
		zbxMap, err = addLink(zbxMap, client, mapping, hosts, options)
		if err != nil {
			return nil, err
		}
	}

	return zbxMap, nil
}

// CreateMap is used to create the given map.
// The map create parameters can also be exported to a file if a file path is specified.
func CreateMap(client *zabbixgosdk.ZabbixService, m *zabbixgosdk.MapCreateParameters, file string) error {
	res, err := client.Map.Create(m)
	if err != nil {
		return err
	}

	if res == nil || len(res.MapIds) == 0 {
		return fmt.Errorf("an empty response was returned when creating the map")
	}

	if file != "" {
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
	}

	return nil
}
