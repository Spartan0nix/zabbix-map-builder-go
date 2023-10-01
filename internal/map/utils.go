package _map

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
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

// mapExistsAskUser is used to ask the user which operation to perform if a map is already using the given name.
func mapExistsAskUser(name string) int {
	var choice int
	reader := bufio.NewReader(os.Stdin)
	msg := fmt.Sprintf("> A map with the name '%s' already exists\n[0] remove the existing map\n[1] add a random id to the new map (example: '<name>-1234') [default]\n\nYour choice : ", name)

	for {
		fmt.Fprintln(os.Stdout, msg)
		s, _ := reader.ReadString('\n')
		s = strings.TrimSpace(s)

		if s == "0" {
			choice = 0
			break
		} else {
			choice = 1
			break
		}
	}

	return choice
}
