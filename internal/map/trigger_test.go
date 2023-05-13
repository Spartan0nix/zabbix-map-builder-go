package _map

import (
	"testing"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

func TestGetTriggerId(t *testing.T) {
	client := zabbixgosdk.NewZabbixService()
	client.SetUrl(ZABBIX_URL)
	client.SetUser(&zabbixgosdk.ApiUser{
		User: ZABBIX_USER,
		Pwd:  ZABBIX_PWD,
	})

	defer client.Logout()

	err := client.Authenticate()
	if err != nil {
		t.Fatalf("error during Zabbix API authentification.\nReason : %v", err)
	}

	h, err := client.Host.Get(&zabbixgosdk.HostGetParameters{
		Filter: map[string]string{
			"host": "Zabbix server",
		},
	})

	if err != nil {
		t.Fatalf("error while retrieving host 'Zabbix server'.\nReason : %v", err)
	}

	if len(h) == 0 {
		t.Fatal("an empty list of hosts was returned")
	}

	triggerId, err := getTriggerId(client, h[0].HostId, "High CPU utilization")
	if err != nil {
		t.Fatalf("error when executing getTriggerId function.\nReason : %v", err)
	}

	if triggerId == "" {
		t.Fatalf("no trigger matching the pattern was found for host '%s'", h[0].HostId)
	}
}

func TestAddLink(t *testing.T) {
	zbxMap := &zabbixgosdk.MapCreateParameters{}
	params := linkParameters{
		localElement:     "1",
		localTrigger:     "11",
		remoteElement:    "2",
		remoteTrigger:    "12",
		linkColor:        "000000",
		triggerLinkColor: "DD0000",
	}

	zbxMap = addLink(zbxMap, &params)

	if len(zbxMap.Links) != 1 {
		t.Fatalf("wrong number of links set in the map.\nExpected : '1'\nReturned : %d", len(zbxMap.Links))
	}

	link := zbxMap.Links[0]

	if link.SelementId1 != "1" {
		t.Fatalf("wrong element id (1) set.\nExpected : '1'\nReturned : %s", link.SelementId1)
	}

	if link.SelementId2 != "2" {
		t.Fatalf("wrong element id (2) set.\nExpected : '2'\nReturned : %s", link.SelementId2)
	}

	if link.Color != "000000" {
		t.Fatalf("wrong color set.\nExpected : '000000'\nReturned : %s", link.Color)
	}

	if link.LinkTriggers[0].TriggerId != "11" {
		t.Fatalf("wrong trigger id (1) set.\nExpected : '11'\nReturned : %s", link.LinkTriggers[0].TriggerId)
	}

	if link.LinkTriggers[0].Color != "DD0000" {
		t.Fatalf("wrong trigger color (1) set.\nExpected : 'DD0000'\nReturned : %s", link.LinkTriggers[0].Color)
	}

	if link.LinkTriggers[1].TriggerId != "12" {
		t.Fatalf("wrong trigger id (2) set.\nExpected : '12'\nReturned : %s", link.LinkTriggers[1].TriggerId)
	}

	if link.LinkTriggers[1].Color != "DD0000" {
		t.Fatalf("wrong trigger color (2) set.\nExpected : 'DD0000'\nReturned : %s", link.LinkTriggers[1].Color)
	}
}
