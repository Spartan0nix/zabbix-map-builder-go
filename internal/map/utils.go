package _map

import (
	"encoding/hex"
	"fmt"
	"strings"

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
	b := strings.HasPrefix(h, "#")
	if b {
		return fmt.Errorf("hexadecimal color should not start with a '#', value parsed '%s'", h)
	}

	if len(h) != 6 {
		return fmt.Errorf("hexadecimal color are coded using 6 hexadecimal caracters, the given value length is %d", len(h))
	}

	_, err := hex.DecodeString(h)
	return err
}
