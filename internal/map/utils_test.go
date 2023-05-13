package _map

import (
	"testing"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

func TestElementExist(t *testing.T) {
	elements := make([]*zabbixgosdk.MapElement, 0)
	elements = append(elements, &zabbixgosdk.MapElement{
		Id:        "1",
		IconIdOff: "1",
	})
	elements = append(elements, &zabbixgosdk.MapElement{
		Id:        "1",
		IconIdOff: "2",
	})
	elements = append(elements, &zabbixgosdk.MapElement{
		Id:        "2",
		IconIdOff: "3",
	})

	b := elementExist("1", elements)
	if !b {
		t.Fatalf("a true bool should be returned when an element with the same id exist in the list")
	}

	b = elementExist("3", elements)
	if b {
		t.Fatalf("a false bool should be returned when an element with the same id does not exist in the list")
	}
}

func TestValidateHexa(t *testing.T) {
	err := validateHexa("DD00FF")
	if err != nil {
		t.Fatalf("error while executing validateHexa function.\nReason : %v", err)
	}
}

func TestValidateHexaFail(t *testing.T) {
	err := validateHexa("#DD00FF")
	if err == nil {
		t.Fatalf("an error should be returned when the given string start with an '#'")
	}

	err = validateHexa("xn;:")
	if err == nil {
		t.Fatalf("an error should be returned when the given string is not in hexadecimal format")
	}

	err = validateHexa("DD00FF00")
	if err == nil {
		t.Fatalf("an error should be returned when the given string is too long")
	}

}
