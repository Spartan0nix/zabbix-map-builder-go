package app

import (
	"log"
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

var testingClient *zabbixgosdk.ZabbixService

func init() {
	var err error
	testingClient, err = api.InitApi(ZABBIX_URL, ZABBIX_USER, ZABBIX_PWD)
	if err != nil {
		log.Fatalf("error while executing InitApi function.\nReason : %v", err)
	}
}

func TestGetUniqueHosts(t *testing.T) {
	m := make([]*zbxMap.Mapping, 0)
	m = append(m, &zbxMap.Mapping{
		LocalHost:  host,
		RemoteHost: host,
	})

	out, err := getUniqueHosts(testingClient, m)
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

func TestGetUniqueImages(t *testing.T) {
	m := make([]*zbxMap.Mapping, 0)
	m = append(m, &zbxMap.Mapping{
		LocalImage:  "Firewall_(64)",
		RemoteImage: "Switch_(64)",
	})

	out, err := getUniqueImages(testingClient, m)
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
