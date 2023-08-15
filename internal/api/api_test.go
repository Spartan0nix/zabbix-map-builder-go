package api

import (
	"log"
	"testing"

	_ "net/http/pprof"

	zabbixgosdk "github.com/Spartan0nix/zabbix-go-sdk/v2"
)

const (
	ZABBIX_URL  = "http://localhost:4444/api_jsonrpc.php"
	ZABBIX_USER = "Admin"
	ZABBIX_PWD  = "zabbix"
)

var testingClient *zabbixgosdk.ZabbixService

func init() {
	testingClient = zabbixgosdk.NewZabbixService()
	testingClient.Auth.Client.Url = ZABBIX_URL
	testingClient.Map.Client.Url = ZABBIX_URL
	testingClient.Trigger.Client.Url = ZABBIX_URL

	err := authenticate(testingClient, ZABBIX_USER, ZABBIX_PWD)
	if err != nil {
		log.Fatalf("error while initializing testing client\nReason : %v", err)
	}
}

func TestInitService(t *testing.T) {
	c, err := initService(ZABBIX_URL)
	if err != nil {
		t.Fatalf("error while executing initService function.\nReason : %v", err)
	}

	defer c.Logout()

	if c == nil {
		t.Fatalf("an nil pointer was returned instead of *zabbixgosdk.ZabbixService")
	}

	if c.Auth.Client.Url == "" {
		t.Fatalf("the Auth client URL was not set correctly")
	}

	if c.Map.Client.Url == "" {
		t.Fatalf("the Map client URL was not set correctly")
	}

	if c.Trigger.Client.Url == "" {
		t.Fatalf("the Trigger client URL was not set correctly")
	}
}

func TestInitServiceFailConnectivity(t *testing.T) {
	c, err := initService("http://localhost:1234/api_jsonrpc.php")
	if err == nil {
		t.Fatalf("an error should be returned when the server is unreachable")
	}

	if c != nil {
		t.Fatalf("a nil pointer should be returned when the server is unreachable instead of *zabbixgosdk.ZabbixService")
	}
}

func TestAuthenticate(t *testing.T) {
	c, err := initService(ZABBIX_URL)
	if err != nil {
		t.Fatalf("error while executing initService function.\nReason : %v", err)
	}

	defer c.Logout()

	err = authenticate(c, ZABBIX_USER, ZABBIX_PWD)
	if err != nil {
		t.Fatalf("error while executing authenticate function.\nReason : %v", err)
	}

	if c.Map.Client.Token == "" {
		t.Fatalf("the Map client Token was not set correctly")
	}

	if c.Trigger.Client.Token == "" {
		t.Fatalf("the Trigger client Token was not set correctly")
	}
}

func TestAuthenticateFail(t *testing.T) {
	c, err := initService("http://localhost:4444/api_jsonrpc.php")
	if err != nil {
		t.Fatalf("error while executing initService function.\nReason : %v", err)
	}

	defer c.Logout()

	err = authenticate(c, "random-user", "random-password")
	if err == nil {
		t.Fatalf("error should be returned when the authentification failed")
	}
}

func TestInitApi(t *testing.T) {
	c, err := InitApi(ZABBIX_URL, ZABBIX_USER, ZABBIX_PWD)
	if err != nil {
		t.Fatalf("error while executing InitApi function.\nReason : %v", err)
	}

	defer c.Logout()

	if c == nil {
		t.Fatalf("an nil pointer was returned instead of *zabbixgosdk.ZabbixService")
	}

	if c.Auth.Client.Url == "" {
		t.Fatalf("the Auth client URL was not set correctly")
	}

	if c.Map.Client.Url == "" {
		t.Fatalf("the Map client URL was not set correctly")
	}

	if c.Trigger.Client.Url == "" {
		t.Fatalf("the Trigger client URL was not set correctly")
	}

	if c.Map.Client.Token == "" {
		t.Fatalf("the Map client Token was not set correctly")
	}

	if c.Trigger.Client.Token == "" {
		t.Fatalf("the Trigger client Token was not set correctly")
	}
}

func TestInitApiFailConnectivity(t *testing.T) {
	c, err := InitApi("http://localhost:1234/api_jsonrpc.php", ZABBIX_USER, ZABBIX_PWD)
	if err == nil {
		t.Fatalf("an error should be returned when the server is unreachable")
	}

	if c != nil {
		t.Fatalf("a nil pointer should be returned when the server is unreachable instead of *zabbixgosdk.ZabbixService")
	}
}

func TestInitApiFailAuth(t *testing.T) {
	c, err := InitApi(ZABBIX_URL, "random-user", "random-password")
	if err == nil {
		t.Fatalf("an error should be returned when the authentification failed")
	}

	if c != nil {
		t.Fatalf("a nil pointer should be returned when the server is unreachable instead of *zabbixgosdk.ZabbixService")
	}
}

func TestLogout(t *testing.T) {
	c, err := InitApi(ZABBIX_URL, ZABBIX_USER, ZABBIX_PWD)
	if err != nil {
		t.Fatalf("error while executing InitApi function.\nReason : %v", err)
	}

	err = Logout(c)
	if err != nil {
		t.Fatalf("error while executing Logout function.\nReason : %v", err)
	}
}

func TestGetHostsId(t *testing.T) {
	hosts, err := GetHostsId(testingClient, map[string]string{
		"Zabbix server": "",
	})
	if err != nil {
		t.Fatalf("error while executing GetHostsId function.\nReason : %v", err)
	}

	if len(hosts) == 0 {
		t.Fatal("an empty map was returned")
	}

	if hosts["Zabbix server"] == "" {
		t.Fatal("no hostid was associated to the 'Zabbix server' host")
	}
}

func TestGetImagesId(t *testing.T) {
	images, err := GetImagesId(testingClient, map[string]string{
		"Cloud_(24)": "",
	})
	if err != nil {
		t.Fatalf("error while executing GetImagesId function.\nReason : %v", err)
	}

	if len(images) == 0 {
		t.Fatal("an empty map was returned")
	}

	if images["Cloud_(24)"] == "" {
		t.Fatal("no imageid was associated to the 'Cloud_(24)' image")
	}
}

func BenchmarkInitService(b *testing.B) {
	for i := 0; i < b.N; i++ {
		initService(ZABBIX_URL)
	}
}

func BenchmarkAuthenticate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		authenticate(testingClient, ZABBIX_USER, ZABBIX_PWD)
	}
}

func BenchmarkInitApi(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InitApi(ZABBIX_URL, ZABBIX_USER, ZABBIX_PWD)
	}
}

func BenchmarkGetHostsId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetHostsId(testingClient, map[string]string{
			"router-1": "",
			"router-2": "",
		})
	}
}

func BenchmarkGetImagesId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetHostsId(testingClient, map[string]string{
			"Switch_(64)":   "",
			"Firewall_(64)": "",
		})
	}
}
