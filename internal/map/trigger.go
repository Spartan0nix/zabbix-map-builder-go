package _map

import (
	"fmt"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

// getTriggerId is used to retrive the triggerId for a given host with a specific pattern (used to filtrer the description field).
func getTriggerId(client *zabbixgosdk.ZabbixService, hostId string, pattern string) (string, error) {
	t, err := client.Trigger.Get(&zabbixgosdk.TriggerGetParameters{
		Output: []string{
			"triggerid",
		},
		HostIds: []string{
			hostId,
		},
		Filter: map[string]string{
			"description": pattern,
		},
	})

	if err != nil {
		return "", err
	}

	if len(t) > 1 {
		return "", fmt.Errorf("more than one trigger was found for the host '%s' with the given pattern '%s'", hostId, pattern)
	}

	return t[0].Id, nil
}

// linkParameters define the parameters required to create a map link between two hosts.s
type linkParameters struct {
	localElement     string
	localTrigger     string
	remoteElement    string
	remoteTrigger    string
	linkColor        string
	triggerLinkColor string
}

// addLink is used to a link between a remote and local hosts for a given mapping.
func addLink(zbxMap *zabbixgosdk.MapCreateParameters, p *linkParameters) *zabbixgosdk.MapCreateParameters {
	link := zabbixgosdk.MapLink{
		SelementId1: p.localElement,
		SelementId2: p.remoteElement,
		Color:       p.linkColor,
		LinkTriggers: []*zabbixgosdk.MapLinkTrigger{
			{
				TriggerId: p.localTrigger,
				Color:     p.triggerLinkColor,
			},
			{
				TriggerId: p.remoteTrigger,
				Color:     p.triggerLinkColor,
			},
		},
	}

	zbxMap.Links = append(zbxMap.Links, &link)

	return zbxMap
}
