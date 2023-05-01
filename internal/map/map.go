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

// BuildMap is used to build a map with the given mapping.
func BuildMap(client *zabbixgosdk.ZabbixService, mappings []*Mapping, hosts map[string]string) (*zabbixgosdk.MapCreateParameters, error) {
	zbxMap := &zabbixgosdk.MapCreateParameters{}
	zbxMap.Name = "test-map"
	zbxMap.Width = "800"
	zbxMap.Height = "800"
	var err error

	for _, mapping := range mappings {
		zbxMap = addHosts(zbxMap, mapping, hosts)
		zbxMap, err = addLink(zbxMap, client, mapping, hosts)
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
