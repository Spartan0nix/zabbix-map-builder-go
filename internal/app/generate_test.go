package app

import (
	"regexp"
	"testing"

	"github.com/Spartan0nix/zabbix-map-builder-go/internal/snmp"
)

func TestGenerateMapping(t *testing.T) {
	e := make([]*snmp.SnmpCdpEntry, 0)
	e = append(e, &snmp.SnmpCdpEntry{
		Address:        []uint8{00, 00, 00, 01},
		DeviceId:       "remote-router",
		Port:           "eth1",
		Capabilities:   []uint8{00, 00, 00, 02},
		LocalPortIndex: "1",
		LocalPort:      "eth0",
	})

	m := generateMapping(e, "local-router", &mappingOptions{
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

func TestBuildTriggerPattern(t *testing.T) {
	re := regexp.MustCompile(`#pattern-to-replace`)
	opts := mappingOptions{
		TriggerPattern: "Interface #pattern-to-replace is down",
	}

	l, r := buildTriggerPattern(re, &opts, "local", "remote")

	if l != "Interface local is down" {
		t.Fatalf("wrong local trigger format returned\nExpected : %s\nReturned : %s", "Interface local is down", l)
	}

	if r != "Interface remote is down" {
		t.Fatalf("wrong remote trigger format returned\nExpected : %s\nReturned : %s", "Interface remote is down", r)
	}
}

func TestBuildTriggerPatternNoUpdate(t *testing.T) {
	re := regexp.MustCompile(`#pattern-to-replace`)
	opts := mappingOptions{
		TriggerPattern: "Interface is down",
	}

	l, r := buildTriggerPattern(re, &opts, "local", "remote")

	if l != "Interface is down" {
		t.Fatalf("wrong local trigger format returned\nExpected : %s\nReturned : %s", "Interface is down", l)
	}

	if r != "Interface is down" {
		t.Fatalf("wrong remote trigger format returned\nExpected : %s\nReturned : %s", "Interface is down", r)
	}
}

func TestReplaceInterface(t *testing.T) {
	re := regexp.MustCompile(`#pattern-to-replace`)
	s := replaceInterface(re, "value : #pattern-to-replace", "new-value")

	if s != "value : new-value" {
		t.Fatalf("wrong value returned\nExpected : %s\nReturned : %s", "value : new-value", s)
	}
}
