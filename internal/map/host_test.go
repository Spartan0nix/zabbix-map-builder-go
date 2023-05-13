package _map

import (
	"testing"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

func TestCreateHostElement(t *testing.T) {
	element := createHostElement("2", "1", "11")
	elementHosts := element.Elements.([]zabbixgosdk.MapElementHost)

	if element.Id != "2" {
		t.Fatalf("wrong element id set.\nExpected : '2'\nReturned : %s", element.Id)
	}

	if len(elementHosts) != 1 {
		t.Fatalf("wrong number of host elements set.\nExpected : 1\nReturned : %d", len(elementHosts))
	}

	if elementHosts[0].Id != "1" {
		t.Fatalf("wrong id set for the host element.\nExpected : '1'\nReturned : %s", elementHosts[0].Id)
	}

	if element.ElementType != zabbixgosdk.MapHost {
		t.Fatalf("wrong type of elementType set.\nExpected : %s\nReturned : %v", zabbixgosdk.MapHost, element.ElementType)
	}

	if element.IconIdOff != "11" {
		t.Fatalf("wrong icon id set.\nExpected : 11\nReturned : %s", element.IconIdOff)
	}
}

func TestAddHosts(t *testing.T) {
	zbxMap := &zabbixgosdk.MapCreateParameters{}
	zbxMap = addHosts(zbxMap, "2", "1", "11")

	if len(zbxMap.Elements) != 1 {
		t.Fatalf("wrong number of elements set.\nExpected : 1\nReturned : %d", len(zbxMap.Elements))
	}
}
