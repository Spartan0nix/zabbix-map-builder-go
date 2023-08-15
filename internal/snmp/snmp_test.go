package snmp

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"

	"github.com/Spartan0nix/zabbix-map-builder-go/internal/logging"
	"github.com/gosnmp/gosnmp"
)

func convertStringToUint(s string) uint8 {
	u, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		log.Fatalf("error while converting string '%s' to uint\nReason : %v", s, err)
	}

	return uint8(u)
}

func TestSupportedCdpCacheCapabilities(t *testing.T) {
	supported := [][]uint8{
		{
			0,
			0,
			0,
			convertStringToUint("28"),
		},
		{
			0,
			0,
			0,
			convertStringToUint("29"),
		},
	}

	for _, value := range supported {
		if ok := SupportedCdpCacheCapabilities(value); !ok {
			t.Fatalf("'%s' should be supported", fmt.Sprintf("%x", value))
		}
	}
}

func TestSupportedCdpCacheCapabilitiesUnsupported(t *testing.T) {
	unsupported := [][]uint8{
		{
			0,
			0,
			0,
			2,
		},
		{
			0,
			0,
			0,
			1,
		},
	}

	for _, value := range unsupported {
		if ok := SupportedCdpCacheCapabilities(value); ok {
			t.Fatalf("'%s' should not be supported", fmt.Sprintf("%x", value[len(value)-1]))
		}
	}
}

func BenchmarkSupportedCdpCacheCapabilities(b *testing.B) {
	supported := []uint8{
		0,
		0,
		0,
		convertStringToUint("28"),
	}

	for i := 0; i < b.N; i++ {
		SupportedCdpCacheCapabilities(supported)
	}
}

func TestParseCdpCache(t *testing.T) {
	expected := SnmpCdpEntry{
		Address: []uint8{
			convertStringToUint("AC"),
			convertStringToUint("10"),
			convertStringToUint("50"),
			convertStringToUint("A2"),
		},
		Capabilities: []uint8{
			0,
			0,
			0,
			convertStringToUint("29"),
		},
		DeviceId:       "router-2",
		LocalPortIndex: "1",
		Port:           "eth0",
	}

	entries := []*SnmpEntry{
		{
			Oid:    ".1.3.6.1.4.1.9.9.23.1.2.1.1.4.1.1",
			Value:  expected.Address,
			Format: gosnmp.OctetString,
		},
		{
			Oid:    ".1.3.6.1.4.1.9.9.23.1.2.1.1.6.1.1",
			Value:  []byte(expected.DeviceId),
			Format: gosnmp.OctetString,
		},
		{
			Oid:    ".1.3.6.1.4.1.9.9.23.1.2.1.1.7.1.1",
			Value:  []byte(expected.Port),
			Format: gosnmp.OctetString,
		},
		{
			Oid:    ".1.3.6.1.4.1.9.9.23.1.2.1.1.9.1.1",
			Value:  expected.Capabilities,
			Format: gosnmp.OctetString,
		},
	}

	res := ParseCdpCache(entries, logging.NewLogger(logging.Warning))

	if len(res) != 1 {
		t.Fatalf("wrong number of cdp entries returned\nExpected : %d\nReturned : %d", 1, len(res))
	}

	if !reflect.DeepEqual(res[0].Address, expected.Address) {
		t.Fatalf("wrong address assignment\nExpected : %v\nReturned : %v", expected.Address, res[0].Address)
	}

	if !reflect.DeepEqual(res[0].Capabilities, expected.Capabilities) {
		t.Fatalf("wrong capabilities assignment\nExpected : %v\nReturned : %v", expected.Capabilities, res[0].Capabilities)
	}

	if res[0].DeviceId != expected.DeviceId {
		t.Fatalf("wrong device id assignment\nExpected : %s\nReturned : %s", expected.DeviceId, res[0].DeviceId)
	}

	if res[0].LocalPortIndex != expected.LocalPortIndex {
		t.Fatalf("wrong local port index assignment\nExpected : %s\nReturned : %s", expected.LocalPortIndex, res[0].LocalPortIndex)
	}

	if res[0].Port != expected.Port {
		t.Fatalf("wrong port assignment\nExpected : %s\nReturned : %s", expected.Port, res[0].LocalPort)
	}
}

func BenchmarkParseCdpCache(b *testing.B) {
	cdpCache := []*SnmpEntry{
		{
			Oid: ".1.3.6.1.4.1.9.9.23.1.2.1.1.4.1.1",
			Value: []uint8{
				convertStringToUint("AC"),
				convertStringToUint("10"),
				convertStringToUint("50"),
				convertStringToUint("A2"),
			},
		},
		{
			Oid:   ".1.3.6.1.4.1.9.9.23.1.2.1.1.6.1.1",
			Value: []byte("router-1"),
		},
		{
			Oid:   ".1.3.6.1.4.1.9.9.23.1.2.1.1.7.1.1",
			Value: []byte("eth0"),
		},
		{
			Oid: ".1.3.6.1.4.1.9.9.23.1.2.1.1.9.1.1",
			Value: []uint8{
				0,
				0,
				0,
				convertStringToUint("29"),
			},
		},
	}

	for i := 0; i < b.N; i++ {
		ParseCdpCache(cdpCache, nil)
	}
}

func TestGetLocalInterfacesName(t *testing.T) {
	connectSnmpClient(snmpClient)
	defer snmpClient.Conn.Close()

	entry := SnmpCdpEntry{
		DeviceId:       "router-1",
		LocalPortIndex: "1",
		Port:           "eth0",
	}

	err := GetLocalInterfacesName(snmpClient, []*SnmpCdpEntry{&entry}, nil)
	if err != nil {
		t.Fatalf("error while executing GetLocalInterfacesName function\nReason : %v", err)
	}

	if entry.LocalPort != "eth0" {
		t.Fatalf("wrong local port assignment\nExpected : %s\nReturned : %s", "eth0", entry.LocalPort)
	}
}

func BenchmarkGetLocalInterfacesName(b *testing.B) {
	connectSnmpClient(snmpClient)
	defer snmpClient.Conn.Close()

	entries := []*SnmpCdpEntry{
		{
			DeviceId:       "router-1",
			LocalPortIndex: "1",
			Port:           "eth0",
		},
	}

	for i := 0; i < b.N; i++ {
		err := GetLocalInterfacesName(snmpClient, entries, nil)
		if err != nil {
			b.Fatalf("error while running benchmark GetLocalInterfacesName\nReason: %v", err)
		}
	}
}

func TestGetHostname(t *testing.T) {
	connectSnmpClient(snmpClient)
	defer snmpClient.Conn.Close()

	h, err := GetHostname(snmpClient, nil)
	if err != nil {
		t.Fatalf("error while executing GetHostname function\nReason : %v", err)
	}

	if h != "router-1" {
		t.Fatalf("wrong hostname returned\nExpected : %s\nReturned : %s", "router-1", h)
	}
}

func BenchmarkGetHostname(b *testing.B) {
	connectSnmpClient(snmpClient)
	defer snmpClient.Conn.Close()

	for i := 0; i < b.N; i++ {
		_, err := GetHostname(snmpClient, nil)
		if err != nil {
			b.Fatalf("error while running benchmark GetHostname\nReason: %v", err)
		}
	}
}
