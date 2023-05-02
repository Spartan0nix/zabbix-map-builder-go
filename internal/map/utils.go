package _map

import (
	"strconv"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

// elementExist is used to check if an element with the given id exist in a list of element.
func elementExist(id string, elements []*zabbixgosdk.MapElement) bool {
	for _, element := range elements {
		if element.Id == id {
			return true
		}
	}

	return false
}

// validateHexa is used to validate that the given string is in hexadecimal format.
func validateHexa(h string) error {
	if string(h[0]) == "#" {
		h = h[:1]
	}

	_, err := strconv.ParseUint(h, 16, 64)
	if err != nil {
		return err
	}

	return nil
}
