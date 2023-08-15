package snmp

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/gosnmp/gosnmp"
)

const (
	router    = "172.16.81.161"
	port      = uint16(1161)
	community = "router-1"
)

var snmpClient *gosnmp.GoSNMP

func init() {
	snmpClient = &gosnmp.GoSNMP{
		Target:    router,
		Port:      port,
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
	}
}

func connectSnmpClient(p *gosnmp.GoSNMP) {
	err := p.Connect()
	if err != nil {
		log.Fatalf("error while initializing the snmp client\nReason : %v", err)
	}
}

func TestWalkBulk(t *testing.T) {
	connectSnmpClient(snmpClient)
	defer snmpClient.Conn.Close()
	oid := "1.3.6.1.2.1.2.2.1.2"

	res, err := WalkBulk(snmpClient, oid)
	if err != nil {
		t.Fatalf("error while executing WalkBulk function\nReason : %v", err)
	}

	if len(res.Entries) != 2 {
		t.Fatalf("wrong number of entries returned\nExpected : %d\nReturned : %d", 2, len(res.Entries))
	}

	for _, r := range res.Entries {
		if !strings.Contains(r.Oid, oid) {
			t.Fatalf("expected the returned oid to contains '%s'\nOid returned : %s", oid, r.Oid)
		}
	}
}

func BenchmarkWalkBulk(b *testing.B) {
	connectSnmpClient(snmpClient)
	defer snmpClient.Conn.Close()
	oid := "1.3.6.1.2.1.2.2.1.2"

	for i := 0; i < b.N; i++ {
		_, err := WalkBulk(snmpClient, oid)
		if err != nil {
			b.Fatalf("error while executing function WalkBulk\nReason : %v", err)
		}
	}
}

func TestGet(t *testing.T) {
	connectSnmpClient(snmpClient)
	defer snmpClient.Conn.Close()

	res, err := Get(snmpClient, []string{
		"1.3.6.1.2.1.2.2.1.2.1",
		"1.3.6.1.2.1.2.2.1.2.2",
	})
	if err != nil {
		t.Fatalf("error while executing Get function\nReason : %v", err)
	}

	if len(res.Entries) != 2 {
		t.Fatalf("wrong number of entries returned\nExpected : %d\nReturned : %d", 2, len(res.Entries))
	}

	for _, r := range res.Entries {
		strValue := string(r.Value.([]uint8))

		switch r.Oid {
		case ".1.3.6.1.2.1.2.2.1.2.1":
			if strValue != "eth0" {
				t.Fatalf("wrong value returned for the following oid : '%s'\nExpected : %s\nReturned : %s", ".1.3.6.1.2.1.2.2.1.2.1", "eth0", strValue)
			}
		case ".1.3.6.1.2.1.2.2.1.2.2":
			if strValue != "eth1" {
				t.Fatalf("wrong value returned for the following oid : '%s'\nExpected : %s\nReturned : %s", ".1.3.6.1.2.1.2.2.1.2.2", "eth1", strValue)
			}
		default:
			t.Fatalf("unexpected value returned\nOid : %s\nValue : %s", r.Oid, strValue)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	connectSnmpClient(snmpClient)
	defer snmpClient.Conn.Close()
	oids := []string{
		"1.3.6.1.2.1.2.2.1.2.1",
		"1.3.6.1.2.1.2.2.1.2.2",
	}

	for i := 0; i < b.N; i++ {
		_, err := Get(snmpClient, oids)
		if err != nil {
			b.Fatalf("error while executing Get function\nReason : %v", err)
		}
	}
}

func TestGetNext(t *testing.T) {
	connectSnmpClient(snmpClient)
	defer snmpClient.Conn.Close()

	res, err := GetNext(snmpClient, []string{
		"1.3.6.1.2.1.2.2.1.2.1",
	})
	if err != nil {
		t.Fatalf("error while executing Get function\nReason : %v", err)
	}

	if len(res.Entries) != 1 {
		t.Fatalf("wrong number of entries returned\nExpected : %d\nReturned : %d", 1, len(res.Entries))
	}

	entry := res.Entries[0]
	if entry.Oid != ".1.3.6.1.2.1.2.2.1.2.2" {
		t.Fatalf("wrong oid returned\nExpected : %s\nReturned : %s", ".1.3.6.1.2.1.2.2.1.2.2", entry.Oid)
	}

	strValue := string(entry.Value.([]uint8))
	if strValue != "eth1" {
		t.Fatalf("wrong value returned\nExpected : %s\nReturned : %s", "eth1", strValue)
	}
}

func BenchmarkGetNext(b *testing.B) {
	connectSnmpClient(snmpClient)
	defer snmpClient.Conn.Close()
	oids := []string{
		"1.3.6.1.2.1.2.2.1.2.1",
	}

	for i := 0; i < b.N; i++ {
		_, err := GetNext(snmpClient, oids)
		if err != nil {
			b.Fatalf("error while executing GetNext function\nReason : %v", err)
		}
	}
}
