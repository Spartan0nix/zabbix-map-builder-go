package _map

import (
	"os"
	"testing"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

func TestElementExist(t *testing.T) {
	elements := []*zabbixgosdk.MapElement{
		{
			Id:        "1",
			IconIdOff: "1",
		},
		{
			Id:        "1",
			IconIdOff: "2",
		},
		{
			Id:        "2",
			IconIdOff: "3",
		},
	}

	b := elementExist("1", elements)
	if !b {
		t.Fatalf("a true bool should be returned when an element with the same id exist in the list")
	}

	b = elementExist("3", elements)
	if b {
		t.Fatalf("a false bool should be returned when an element with the same id does not exist in the list")
	}
}

func BenchmarkElementExist(b *testing.B) {
	elements := []*zabbixgosdk.MapElement{
		{
			Id:        "1",
			IconIdOff: "1",
		},
		{
			Id:        "1",
			IconIdOff: "2",
		},
		{
			Id:        "2",
			IconIdOff: "3",
		},
	}

	for i := 0; i < b.N; i++ {
		elementExist("1", elements)
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

func BenchmarkValidateHexa(b *testing.B) {
	hexa := "FFFFFF"
	for i := 0; i < b.N; i++ {
		validateHexa(hexa)
	}
}

func TestMapExistsAskUser(t *testing.T) {
	_, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	res := mapExistsAskUser("test-map-builder")
	if res != 1 {
		t.Fatalf("wrong value returned\nExpected : 1\nReturned : %d", res)
	}

	os.Stdout = oldStdout
}

func BenchmarkMapExistsAskUser(b *testing.B) {
	_, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	for i := 0; i < b.N; i++ {
		mapExistsAskUser("test-map-builder")
	}

	os.Stdout = oldStdout
}
