package _map

import (
	"fmt"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

// getTriggerId is used to retrive the triggerId for a given host with a specific pattern (used to filtrer the description field).
func getTriggerId(client *zabbixgosdk.ZabbixService, id string, pattern string) (string, error) {
	t, err := client.Trigger.Get(&zabbixgosdk.TriggerGetParameters{
		Output: []string{
			"triggerid",
		},
		HostIds: []string{
			id,
		},
		Filter: map[string]string{
			"description": pattern,
		},
	})

	if err != nil {
		return "", err
	}

	if len(t) > 1 {
		return "", fmt.Errorf("more than one trigger was found for the host '%s' with the given pattern '%s'", id, pattern)
	}

	return t[0].Id, nil
}

// addLink is used to a link between a remote and local hosts for a given mapping.
func addLink(zbxMap *zabbixgosdk.MapCreateParameters, client *zabbixgosdk.ZabbixService, mapping *Mapping, hosts map[string]string) (*zabbixgosdk.MapCreateParameters, error) {
	localTriggerId, err := getTriggerId(client, hosts[mapping.LocalHost], mapping.LocalTriggerPattern)
	if err != nil {
		return nil, err
	}

	remoteTriggerId, err := getTriggerId(client, hosts[mapping.RemoteHost], mapping.RemoteTriggerPattern)
	if err != nil {
		return nil, err
	}

	link := zabbixgosdk.MapLink{
		SelementId1: hosts[mapping.LocalHost],
		SelementId2: hosts[mapping.RemoteHost],
		LinkTriggers: []*zabbixgosdk.MapLinkTrigger{
			{
				TriggerId: localTriggerId,
			},
			{
				TriggerId: remoteTriggerId,
			},
		},
	}

	zbxMap.Links = append(zbxMap.Links, &link)

	return zbxMap, nil
}
