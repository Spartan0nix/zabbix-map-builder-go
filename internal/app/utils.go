package app

import (
	"fmt"
	"os"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
	"github.com/Spartan0nix/zabbix-map-builder-go/internal/api"
	zbxMap "github.com/Spartan0nix/zabbix-map-builder-go/internal/map"
)

// outputToFile is used to write data to a file.
func outputToFile(file string, b []byte) error {
	if file == "" {
		return fmt.Errorf("file name cannot be empty")
	}

	err := os.WriteFile(file, b, 0640)
	return err
}

// getUniqueHosts is used to get a map where each key correspond to an host name reference in the list of Mapping and the value, the hostid associated on the Zabbix server.
func getUniqueHosts(client *zabbixgosdk.ZabbixService, mappings []*zbxMap.Mapping) (map[string]string, error) {
	out := make(map[string]string, 0)

	for _, m := range mappings {
		_, exist := out[m.LocalHost]
		if !exist {
			out[m.LocalHost] = ""
		}

		_, exist = out[m.RemoteHost]
		if !exist {
			out[m.RemoteHost] = ""
		}
	}

	// Retrieve the id of each hosts and provide a mapping 'host' -> 'hostid'
	out, err := api.GetHostsId(client, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// getUniqueHosts is used to get a map where each key correspond to an image name reference in the list of Mapping and the value, the imageid associated on the Zabbix server.
func getUniqueImages(client *zabbixgosdk.ZabbixService, mappings []*zbxMap.Mapping) (map[string]string, error) {
	out := make(map[string]string, 0)

	for _, m := range mappings {
		_, exist := out[m.LocalImage]
		if !exist {
			out[m.LocalImage] = ""
		}

		_, exist = out[m.RemoteImage]
		if !exist {
			out[m.RemoteImage] = ""
		}
	}

	// Retrieve the id of each images and provide a mapping 'image' -> 'imageid'
	out, err := api.GetImagesId(client, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
