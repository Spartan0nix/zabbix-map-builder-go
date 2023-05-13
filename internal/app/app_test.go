package app

import (
	"os"
	"testing"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

func TestOutputToFile(t *testing.T) {
	file := "test-output-file.json"

	err := outputToFile(file, &zabbixgosdk.MapCreateParameters{
		Map: zabbixgosdk.Map{
			Name: "test-zabbix-map-builder",
		},
	})

	if err != nil {
		t.Fatalf("error while execution outputToFile function.\nReason : %v", err)
	}

	info, err := os.Stat(file)
	if err != nil {
		t.Fatalf("error while retrieving information about the file '%s'.\nReason : %v", file, err)
	}

	if info == nil {
		t.Fatalf("no info was retrieved about the file '%s'", file)
	}

	err = os.Remove(file)
	if err != nil {
		t.Fatalf("error while removing file '%s'.\nReason : %v", file, err)
	}
}

func TestOutputToFileEmptyName(t *testing.T) {
	err := outputToFile("", &zabbixgosdk.MapCreateParameters{
		Map: zabbixgosdk.Map{
			Name: "test-zabbix-map-builder",
		},
	})

	if err == nil {
		t.Fatal("an error should be returned when an empty file name is passed to the outputToFile function")
	}
}
