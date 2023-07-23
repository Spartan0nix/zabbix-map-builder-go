package app

import (
	"regexp"
	"testing"

	"github.com/Spartan0nix/zabbix-map-builder-go/internal/snmp"
)

var SnmpCdpEntries = []*snmp.SnmpCdpEntry{
	{
		Address:        []uint8{00, 00, 00, 01},
		DeviceId:       "remote-router",
		Port:           "eth1",
		Capabilities:   []uint8{00, 00, 00, 02},
		LocalPortIndex: "1",
		LocalPort:      "eth0",
	},
}

var re = regexp.MustCompile(`#pattern-to-replace`)
var reString = `#pattern-to-replace`

func TestGenerateMapping(t *testing.T) {
	m := generateMapping(SnmpCdpEntries, "local-router", &mappingOptions{
		TriggerPattern: "Interface #INTERFACE is down",
		LocalImage:     "Switch_(64)",
		RemoteImage:    "Firewall_(64)",
	})

	if len(m) != 1 {
		t.Fatalf("wrong length returned\nExpected : %d\nReturned : %d", 1, len(m))
	}

	// ---
	if m[0].LocalHost != "local-router" {
		t.Fatalf("wrong 'LocalHost' field assignment\nExpected : %s\nReturned : %s", "local-router", m[0].LocalHost)
	}

	if m[0].LocalInterface != "eth0" {
		t.Fatalf("wrong 'LocalInterface' field assignment\nExpected : %s\nReturned : %s", "eth0", m[0].LocalInterface)
	}

	if m[0].LocalTriggerPattern != "Interface eth0 is down" {
		t.Fatalf("wrong 'LocalInterface' field assignment\nExpected : %s\nReturned : %s", "Interface eth0 is down", m[0].LocalTriggerPattern)
	}

	if m[0].LocalImage != "Switch_(64)" {
		t.Fatalf("wrong 'LocalImage' field assignment\nExpected : %s\nReturned : %s", "Switch_(64)", m[0].LocalImage)
	}

	// ---
	if m[0].RemoteHost != "remote-router" {
		t.Fatalf("wrong 'RemoteHost' field assignment\nExpected : %s\nReturned : %s", "remote-router", m[0].RemoteHost)
	}
	if m[0].RemoteInterface != "eth1" {
		t.Fatalf("wrong 'RemoteInterface' field assignment\nExpected : %s\nReturned : %s", "eth1", m[0].RemoteInterface)
	}

	if m[0].RemoteTriggerPattern != "Interface eth1 is down" {
		t.Fatalf("wrong 'RemoteTriggerPattern' field assignment\nExpected : %s\nReturned : %s", "Interface eth1 is down", m[0].RemoteTriggerPattern)
	}

	if m[0].RemoteImage != "Firewall_(64)" {
		t.Fatalf("wrong 'RemoteImage' field assignment\nExpected : %s\nReturned : %s", "Firewall_(64)", m[0].RemoteImage)
	}
}

func BenchmarkGenerateMapping(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateMapping(SnmpCdpEntries, "local-router", &mappingOptions{
			TriggerPattern: "",
			LocalImage:     "Switch_(64)",
			RemoteImage:    "Firewall_(64)",
		})
	}
}

func BenchmarkGenerateMappingTriggerPattern(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateMapping(SnmpCdpEntries, "local-router", &mappingOptions{
			TriggerPattern: "Interface #INTERFACE is down",
			LocalImage:     "Switch_(64)",
			RemoteImage:    "Firewall_(64)",
		})
	}
}

func TestBuildTriggerPattern(t *testing.T) {
	opts := mappingOptions{
		TriggerPattern: "Interface #pattern-to-replace is down",
	}

	l, r := buildTriggerPattern(re, reString, &opts, "local", "remote")

	if l != "Interface local is down" {
		t.Fatalf("wrong local trigger format returned\nExpected : %s\nReturned : %s", "Interface local is down", l)
	}

	if r != "Interface remote is down" {
		t.Fatalf("wrong remote trigger format returned\nExpected : %s\nReturned : %s", "Interface remote is down", r)
	}
}

func TestBuildTriggerPatternNoUpdate(t *testing.T) {
	opts := mappingOptions{
		TriggerPattern: "Interface is down",
	}

	l, r := buildTriggerPattern(re, reString, &opts, "local", "remote")

	if l != "Interface is down" {
		t.Fatalf("wrong local trigger format returned\nExpected : %s\nReturned : %s", "Interface is down", l)
	}

	if r != "Interface is down" {
		t.Fatalf("wrong remote trigger format returned\nExpected : %s\nReturned : %s", "Interface is down", r)
	}
}

func BenchmarkBuildTriggerPattern(b *testing.B) {
	opts := mappingOptions{
		TriggerPattern: "Interface #INTERFACE is down",
	}

	for i := 0; i < b.N; i++ {
		buildTriggerPattern(re, reString, &opts, "local", "remote")
	}
}

func BenchmarkBuildTriggerPatternNoUpdate(b *testing.B) {
	opts := mappingOptions{
		TriggerPattern: "Interface is down",
	}

	for i := 0; i < b.N; i++ {
		buildTriggerPattern(re, reString, &opts, "local", "remote")
	}
}

func TestReplaceInterface(t *testing.T) {
	re := regexp.MustCompile(`#pattern-to-replace`)
	s := replaceInterface(re, `#pattern-to-replace`, "value : #pattern-to-replace", "new-value")

	if s != "value : new-value" {
		t.Fatalf("wrong value returned\nExpected : %s\nReturned : %s", "value : new-value", s)
	}
}

func BenchmarkReplaceInterface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		replaceInterface(re, reString, "Interface #INTERFACE is down", "eth0")
	}
}

func BenchmarkReplaceInterfaceNoUpdate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		replaceInterface(re, reString, "Interface is down", "eth0")
	}
}
