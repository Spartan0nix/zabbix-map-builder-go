package _map

import zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"

// elementExist is used to check if an element with the given id exist in a list of element.
func elementExist(id string, elements []*zabbixgosdk.MapElement) bool {
	for _, element := range elements {
		if element.Id == id {
			return true
		}
	}

	return false
}
