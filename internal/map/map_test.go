package _map

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

const (
	ZABBIX_URL  = "http://localhost:4444/api_jsonrpc.php"
	ZABBIX_USER = "Admin"
	ZABBIX_PWD  = "zabbix"
)

// generateMapName is used to generate a random name for each map created during test.
func generateMapName() string {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(rand.Intn(9999))

	return fmt.Sprintf("test-map-builder-%d", value)
}

func TestValidate(t *testing.T) {
	opts := MapOptions{
		Name:         "test-map",
		Color:        "EEDDFF",
		TriggerColor: "EEDDFF",
		Mappings: []*Mapping{
			{
				LocalHost: "local-host",
			},
		},
		Hosts: map[string]string{
			"local-host": "",
		},
		Images: map[string]string{
			"Firewall_(64)": "",
		},
	}

	err := opts.Validate()
	if err != nil {
		t.Fatalf("error while executing Validate function.\nReason : %v", err)
	}
}

func TestValidateEmptyValue(t *testing.T) {
	opts := MapOptions{
		Name:         "test-map",
		Color:        "",
		TriggerColor: "",
		Mappings: []*Mapping{
			{
				LocalHost: "local-host",
			},
		},
		Hosts: map[string]string{
			"local-host": "",
		},
		Images: map[string]string{
			"Firewall_(64)": "",
		},
	}

	err := opts.Validate()
	if err != nil {
		t.Fatalf("error while executing Validate function.\nReason : %v", err)
	}

	if opts.Color != "000000" {
		t.Fatalf("wrong default color returned\nExpected : 000000\nReturned : %s", opts.Color)
	}

	if opts.TriggerColor != "DD0000" {
		t.Fatalf("wrong default trigger color returned\nExpected : DD0000\nReturned : %s", opts.TriggerColor)
	}
}

func TestValidateFailName(t *testing.T) {
	opts := MapOptions{
		Name:         "",
		Color:        "000000",
		TriggerColor: "DD0000",
		Mappings: []*Mapping{
			{
				LocalHost: "local-host",
			},
		},
		Hosts: map[string]string{
			"local-host": "",
		},
		Images: map[string]string{
			"Firewall_(64)": "",
		},
	}

	err := opts.Validate()
	if err == nil {
		t.Fatalf("an error should be returned when name checking fails")
	}
}

func TestValidateFailColor(t *testing.T) {
	opts := MapOptions{
		Name:         "test-map",
		Color:        "xm^ù",
		TriggerColor: "DD0000",
		Mappings: []*Mapping{
			{
				LocalHost: "local-host",
			},
		},
		Hosts: map[string]string{
			"local-host": "",
		},
		Images: map[string]string{
			"Firewall_(64)": "",
		},
	}

	err := opts.Validate()
	if err == nil {
		t.Fatalf("an error should be returned when color checking fails")
	}
}

func TestValidateFailTriggerColor(t *testing.T) {
	opts := MapOptions{
		Name:         "test-map",
		Color:        "000000",
		TriggerColor: "xvù*",
		Mappings: []*Mapping{
			{
				LocalHost: "local-host",
			},
		},
		Hosts: map[string]string{
			"local-host": "",
		},
		Images: map[string]string{
			"Firewall_(64)": "",
		},
	}

	err := opts.Validate()
	if err == nil {
		t.Fatalf("an error should be returned when trigger color checking fails")
	}
}

func TestValidateFailMappings(t *testing.T) {
	opts := MapOptions{
		Name:         "test-map",
		Color:        "000000",
		TriggerColor: "DD0000",
		Hosts: map[string]string{
			"local-host": "",
		},
		Images: map[string]string{
			"Firewall_(64)": "",
		},
	}

	err := opts.Validate()
	if err == nil {
		t.Fatalf("an error should be returned when mapping checking fails")
	}
}

func TestValidateFailHosts(t *testing.T) {
	opts := MapOptions{
		Name:         "test-map",
		Color:        "000000",
		TriggerColor: "DD0000",
		Mappings: []*Mapping{
			{
				LocalHost: "local-host",
			},
		},
		Images: map[string]string{
			"Firewall_(64)": "",
		},
	}

	err := opts.Validate()
	if err == nil {
		t.Fatalf("an error should be returned when hosts checking fails")
	}
}

func TestValidateFailImages(t *testing.T) {
	opts := MapOptions{
		Name:         "test-map",
		Color:        "000000",
		TriggerColor: "xvù*",
		Mappings: []*Mapping{
			{
				LocalHost: "local-host",
			},
		},
		Hosts: map[string]string{
			"local-host": "",
		},
	}

	err := opts.Validate()
	if err == nil {
		t.Fatalf("an error should be returned when images checking fails")
	}
}

func TestBuildElementsId(t *testing.T) {
	zbxMap := &zabbixgosdk.MapCreateParameters{
		Elements: []*zabbixgosdk.MapElement{
			{
				Id: "1",
				Elements: []zabbixgosdk.MapElementHost{
					{
						Id: "1",
					},
				},
			},
		},
	}

	local, remote := buildElementsId(zbxMap, "1", "2")

	if local != "1-2" {
		t.Fatalf("local id should have been incremented.\nExpected : '1-1'\nReturned : %s", local)
	}

	if remote != "2" {
		t.Fatalf("remote id should not have been incremented.\nExpected : '2'\nReturned : %s", remote)
	}
}

func TestBuildElementsIdNoIncrement(t *testing.T) {
	zbxMap := &zabbixgosdk.MapCreateParameters{}

	local, remote := buildElementsId(zbxMap, "1", "2")

	if local != "1" {
		t.Fatalf("local id should not have been incremented.\nExpected : '1'\nReturned : %s", local)
	}

	if remote != "2" {
		t.Fatalf("remote id should not have been incremented.\nExpected : '2'\nReturned : %s", remote)
	}
}

func TestCreateMap(t *testing.T) {
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

	err = CreateMap(client, &zabbixgosdk.MapCreateParameters{
		Map: zabbixgosdk.Map{
			Height: "800",
			Width:  "800",
			Name:   generateMapName(),
		},
	})

	if err != nil {
		t.Fatalf("error when executing CreateMap function.\nReason : %v", err)
	}
}
