package _map

import (
	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

// addHostElement is used to add a new MapElementHost to a MapCreateParameters.
// The given id is used to reference the host in the map links.
func addHostElement(zbxMap *zabbixgosdk.MapCreateParameters, id string) *zabbixgosdk.MapElement {
	return &zabbixgosdk.MapElement{
		Id: id,
		Elements: []zabbixgosdk.MapElementHost{
			{
				Id: id,
			},
		},
		ElementType: zabbixgosdk.MapHost,
		// -- TODO : accept dynamic value
		IconIdOff: "156",
	}
}

// addHosts is used to add hosts (local and remote) for a given mapping if they do not already exist in the map.
func addHosts(zbxMap *zabbixgosdk.MapCreateParameters, mapping *Mapping, hosts map[string]string) *zabbixgosdk.MapCreateParameters {
	if exist := elementExist(hosts[mapping.LocalHost], zbxMap.Elements); !exist {
		element := addHostElement(zbxMap, hosts[mapping.LocalHost])
		zbxMap.Elements = append(zbxMap.Elements, element)
	}

	if exist := elementExist(hosts[mapping.RemoteHost], zbxMap.Elements); !exist {
		element := addHostElement(zbxMap, hosts[mapping.RemoteHost])
		zbxMap.Elements = append(zbxMap.Elements, element)
	}

	return zbxMap
}
