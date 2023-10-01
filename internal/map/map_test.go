package _map

import (
	"fmt"
	"log"
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

var testingClient *zabbixgosdk.ZabbixService
var routersId = make(map[string]string, 0)

func init() {
	testingClient = zabbixgosdk.NewZabbixService()
	testingClient.SetUrl(ZABBIX_URL)
	testingClient.SetUser(&zabbixgosdk.ApiUser{
		User: ZABBIX_USER,
		Pwd:  ZABBIX_PWD,
	})

	err := testingClient.Authenticate()
	if err != nil {
		log.Fatalf("error during Zabbix API authentification.\nReason : %v", err)
	}

	// Retrieve the hosts id
	routers, err := testingClient.Host.Get(&zabbixgosdk.HostGetParameters{
		Filter: map[string][]string{
			"host": {
				"router-1",
				"router-2",
				"router-3",
			},
		},
		Output: []string{
			"hostid",
			"name",
		},
	})

	if err != nil {
		log.Fatalf("error while retrieving the hosts id\nReason : %v", err)
	}

	for _, router := range routers {
		routersId[router.Name] = router.HostId
	}
}

// generateMapName is used to generate a random name for each map created during test.
func generateMapName() string {
	rand.NewSource(time.Now().UnixNano())
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
		Color:        "xmù",
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

func BenchmarkValidate(b *testing.B) {
	var opts = MapOptions{
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

	for i := 0; i < b.N; i++ {
		opts.Validate()
	}
}

func TestBuildMap(t *testing.T) {
	opts := &MapOptions{
		Name:         "test-build-map",
		Color:        "000000",
		TriggerColor: "DD0000",
		Height:       600,
		Width:        600,
		Spacer:       100,
		StackHosts:   true,
		Mappings: []*Mapping{
			{
				LocalHost:            "router-1",
				LocalInterface:       "eth0",
				LocalTriggerPattern:  "Interface eth0(): Link down",
				LocalImage:           "Firewall_(64)",
				RemoteHost:           "router-2",
				RemoteInterface:      "eth0",
				RemoteTriggerPattern: "Interface eth0(): Link down",
				RemoteImage:          "Switch_(64)",
			},
		},
		Hosts: map[string]string{
			"router-1": routersId["router-1"],
			"router-2": routersId["router-2"],
		},
		Images: map[string]string{
			"Firewall_(64)": "10",
			"Switch_(64)":   "11",
		},
	}

	zbxMap, err := BuildMap(testingClient, opts, nil)
	if err != nil {
		t.Fatalf("error while executing BuildMap function\nReason : %v", err)
	}

	if zbxMap.Name != "test-build-map" {
		t.Fatalf("wrong map name assigned\nExpected : %s\nReturned : %s", "test-build-map", zbxMap.Name)
	}

	if zbxMap.Height != "600" {
		t.Fatalf("wrong map height assigned\nExpected : %s\nReturned : %s", "600", zbxMap.Height)
	}

	if zbxMap.Width != "600" {
		t.Fatalf("wrong map width assigned\nExpected : %s\nReturned : %s", "600", zbxMap.Width)
	}

	if len(zbxMap.Elements) != 2 {
		t.Fatalf("wrong number of map elements\nExpected : %s\nReturned : %d", "2", len(zbxMap.Elements))
	}
}

func TestBuildMapUnstackedHosts(t *testing.T) {
	opts := &MapOptions{
		Name:         "test-build-map-unstacked",
		Color:        "000000",
		TriggerColor: "DD0000",
		Height:       600,
		Width:        600,
		Spacer:       100,
		StackHosts:   false,
		Mappings: []*Mapping{
			{
				LocalHost:            "router-1",
				LocalInterface:       "eth0",
				LocalTriggerPattern:  "Interface eth0(): Link down",
				LocalImage:           "Firewall_(64)",
				RemoteHost:           "router-2",
				RemoteInterface:      "eth0",
				RemoteTriggerPattern: "Interface eth0(): Link down",
				RemoteImage:          "Switch_(64)",
			},
			{
				LocalHost:            "router-1",
				LocalInterface:       "eth1",
				LocalTriggerPattern:  "Interface eth1(): Link down",
				LocalImage:           "Firewall_(64)",
				RemoteHost:           "router-2",
				RemoteInterface:      "eth1",
				RemoteTriggerPattern: "Interface eth1(): Link down",
				RemoteImage:          "Switch_(64)",
			},
		},
		Hosts: map[string]string{
			"router-1": routersId["router-1"],
			"router-2": routersId["router-2"],
		},
		Images: map[string]string{
			"Firewall_(64)": "10",
			"Switch_(64)":   "11",
		},
	}

	zbxMap, err := BuildMap(testingClient, opts, nil)
	if err != nil {
		t.Fatalf("error while executing BuildMap function\nReason : %v", err)
	}

	if zbxMap.Name != "test-build-map-unstacked" {
		t.Fatalf("wrong map name assigned\nExpected : %s\nReturned : %s", "test-build-map-unstacked", zbxMap.Name)
	}

	if zbxMap.Height != "600" {
		t.Fatalf("wrong map height assigned\nExpected : %s\nReturned : %s", "600", zbxMap.Height)
	}

	if zbxMap.Width != "600" {
		t.Fatalf("wrong map width assigned\nExpected : %s\nReturned : %s", "600", zbxMap.Width)
	}

	for _, r := range zbxMap.Elements {
		t.Log(r.Id)
		t.Log(r.Elements.([]zabbixgosdk.MapElementHost)[0].Id)
		t.Log("---")
	}

	if len(zbxMap.Elements) != 4 {
		t.Fatalf("wrong number of map elements\nExpected : %s\nReturned : %d", "4", len(zbxMap.Elements))
	}
}

func BenchmarkBuildMap(b *testing.B) {
	// Build the mapOptions
	opts := &MapOptions{
		Name:         "test-build-map",
		Color:        "000000",
		TriggerColor: "DD0000",
		Height:       600,
		Width:        600,
		Spacer:       100,
		StackHosts:   true,
		Mappings: []*Mapping{
			{
				LocalHost:            "router-1",
				LocalInterface:       "eth0",
				LocalTriggerPattern:  "Interface eth0(): Link down",
				LocalImage:           "Firewall_(64)",
				RemoteHost:           "router-2",
				RemoteInterface:      "eth0",
				RemoteTriggerPattern: "Interface eth0(): Link down",
				RemoteImage:          "Switch_(64)",
			},
		},
		Hosts: map[string]string{
			"router-1": routersId["router-1"],
			"router-2": routersId["router-2"],
		},
		Images: map[string]string{
			"Firewall_(64)": "10",
			"Switch_(64)":   "11",
		},
	}

	for i := 0; i < b.N; i++ {
		BuildMap(testingClient, opts, nil)
	}
}

func TestCreateMap(t *testing.T) {
	err := CreateMap(testingClient, &zabbixgosdk.MapCreateParameters{
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

func BenchmarkCreateMap(b *testing.B) {
	params := &zabbixgosdk.MapCreateParameters{
		Map: zabbixgosdk.Map{
			Height: "800",
			Width:  "800",
		},
	}

	for i := 0; i < b.N; i++ {
		params.Name = generateMapName()
		CreateMap(testingClient, params)
	}
}
