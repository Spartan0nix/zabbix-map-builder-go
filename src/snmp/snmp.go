package snmp

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
)

func Snmp_init(target string, port uint16, community string) *gosnmp.GoSNMP {
	snmp := &gosnmp.GoSNMP{
		Target:    target,
		Port:      port,
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
	}

	err := snmp.Connect()
	if err != nil {
		log.Fatalf("Error while executing snmp Connect() : %v", err)
	}

	return snmp
}

func Snmp_get_local_hostname(snmp *gosnmp.GoSNMP) string {
	local_hostname_oid := []string{"1.3.6.1.4.1.9.9.23.1.3.4.0"}
	res, err := snmp.Get(local_hostname_oid)
	if err != nil {
		log.Fatalf("Error while retrieve oid '%s' : .Reason : %v", local_hostname_oid[0], err)
	}

	if len(res.Variables) == 0 {
		log.Fatalf("Error while retrieving local_hostname.")
	}

	variable := res.Variables[0]

	return string(variable.Value.([]byte))
}

func Snmp_get_local_interface(snmp *gosnmp.GoSNMP, index string) string {
	local_interface_oid := "1.3.6.1.4.1.9.9.23.1.1.1.1.6." + index
	local_interface_list := []string{local_interface_oid}
	res, err := snmp.Get(local_interface_list)
	if err != nil {
		log.Fatalf("Error while retrieve oid '%s' : .Reason : %v", local_interface_oid, err)
	}

	if len(res.Variables) == 0 {
		log.Fatalf("Error while retrieving local_interface.")
	}

	b := res.Variables[0].Value.([]byte)

	return string(b)
}

func Snmp_get_remote_hostname(snmp *gosnmp.GoSNMP, index string) string {
	remote_hostname_oid := "1.3.6.1.4.1.9.9.23.1.2.1.1.6." + index
	res, err := snmp.BulkWalkAll(remote_hostname_oid)
	if err != nil {
		log.Fatalf("Error while retrieving oid '%s' : .Reason : %v", remote_hostname_oid, err)
	}

	b := res[0].Value.([]byte)
	return string(b)
}

func Snmp_get_remote_interface(snmp *gosnmp.GoSNMP, index string) string {
	remote_interface_oid := "1.3.6.1.4.1.9.9.23.1.2.1.1.7." + index
	res, err := snmp.BulkWalkAll(remote_interface_oid)
	if err != nil {
		log.Fatalf("Error while retrieving oid '%s' : .Reason : %v", remote_interface_oid, err)
	}

	b := res[0].Value.([]byte)
	return string(b)
}

func Extract_ip(pdus []gosnmp.SnmpPDU) []string {
	ips := make([]string, 0)
	for _, pdu := range pdus {
		ip := make([]string, 0)
		bytes := pdu.Value.([]byte)
		for _, b := range bytes {
			str := strconv.FormatUint(uint64(b), 10)
			ip = append(ip, str)
		}
		ips = append(ips, strings.Join(ip, "."))
	}
	return ips
}

func Extract_index(pdus []gosnmp.SnmpPDU) []string {
	indexes := make([]string, 0)
	for _, pdu := range pdus {
		var re = regexp.MustCompile(`(?m)\.[0-9]\.[0-9]\.[0-9]\.[0-9]\.[0-9]\.[0-9]\.[0-9]\.[0-9]\.[0-9]{2}\.[0-9]\.[0-9]\.[0-9]\.[0-9]\.[0-9]\.([0-9]+)\.[0-9]`)

		for _, match := range re.FindAllStringSubmatch(pdu.Name, -1) {
			indexes = append(indexes, match[1])
		}
	}

	return indexes
}
