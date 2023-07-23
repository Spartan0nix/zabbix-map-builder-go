package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
	"github.com/Spartan0nix/zabbix-map-builder-go/internal/api"
	zbxMap "github.com/Spartan0nix/zabbix-map-builder-go/internal/map"
)

const (
	ZABBIX_URL  = "http://localhost:4444/api_jsonrpc.php"
	ZABBIX_USER = "Admin"
	ZABBIX_PWD  = "zabbix"
	host        = "Zabbix server"
)

var hostMapping = []*zbxMap.Mapping{
	{
		LocalHost:  host,
		RemoteHost: host,
	},
}

var imageMapping = []*zbxMap.Mapping{
	{
		LocalImage:  "Firewall_(64)",
		RemoteImage: "Switch_(64)",
	},
}

var testingClient *zabbixgosdk.ZabbixService

func init() {
	var err error
	testingClient, err = api.InitApi(ZABBIX_URL, ZABBIX_USER, ZABBIX_PWD)
	if err != nil {
		log.Fatalf("error while executing InitApi function.\nReason : %v", err)
	}
}

func TestOutputToFile(t *testing.T) {
	file := "test-output-file.json"
	p := &zabbixgosdk.MapCreateParameters{
		Map: zabbixgosdk.Map{
			Name: "test-zabbix-map-builder",
		},
	}

	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("error while marshaling data.\nReason : %v", err)
	}

	err = outputToFile(file, b)
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
	b := make([]byte, 0)

	err := outputToFile("", b)
	if err == nil {
		t.Fatal("an error should be returned when an empty file name is passed to the outputToFile function")
	}
}

func BenchmarkOutputToFile(b *testing.B) {
	data := []byte("random-test")

	for i := 0; i < b.N; i++ {
		fileName := fmt.Sprintf("benchmark-output-file-%d", i)

		outputToFile(fileName, data)
		os.Remove(fileName)
	}
}

func TestGetUniqueHosts(t *testing.T) {
	out, err := getUniqueHosts(testingClient, hostMapping)
	if err != nil {
		t.Fatalf("error while executing getUniqueHosts function.\nReason : %v", err)
	}

	if _, exist := out["Zabbix server"]; !exist {
		t.Fatalf("missing key 'Zabbix server' in the returned map.\nReturned : %v", out)
	}

	if out["Zabbix server"] == "" {
		t.Fatalf("no hostid was associated with the key 'Zabbix server' in the returned map.\nReturned : %v", out)
	}
}

func BenchmarkGetUniqueHosts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getUniqueHosts(testingClient, hostMapping)
	}
}

func TestGetUniqueImages(t *testing.T) {
	out, err := getUniqueImages(testingClient, imageMapping)
	if err != nil {
		t.Fatalf("error while executing getUniqueImages function.\nReason : %v", err)
	}

	if _, exist := out["Firewall_(64)"]; !exist {
		t.Fatalf("missing key 'Firewall_(64)' in the returned map.\nReturned : %v", out)
	}

	if _, exist := out["Switch_(64)"]; !exist {
		t.Fatalf("missing key 'Switch_(64)' in the returned map.\nReturned : %v", out)
	}

	if out["Firewall_(64)"] == "" {
		t.Fatalf("no imageid was associated with the key 'Firewall_(64)' in the returned map.\nReturned : %v", out)
	}

	if out["Switch_(64)"] == "" {
		t.Fatalf("no hostid was associated with the key 'Switch_(64)' in the returned map.\nReturned : %v", out)
	}
}

func BenchmarkGetUniqueImages(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getUniqueImages(testingClient, imageMapping)
	}
}
