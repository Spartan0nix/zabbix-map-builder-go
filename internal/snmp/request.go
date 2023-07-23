package snmp

import (
	"time"

	"github.com/gosnmp/gosnmp"
)

// WalkBulk is used to retrieve subtree data from a oid using a walkbulk request
func WalkBulk(p *gosnmp.GoSNMP, oid string) (*SnmpResponse, error) {
	response := &SnmpResponse{}

	// Keep stats about the duration of the request
	p.OnSent = func(gs *gosnmp.GoSNMP) {
		response.Start = time.Now()
	}
	p.OnRecv = func(gs *gosnmp.GoSNMP) {
		response.End = time.Now()
		response.Duration = response.End.Sub(response.Start)
	}

	// Execute the request
	err := p.BulkWalk(oid, func(dataUnit gosnmp.SnmpPDU) error {
		response.Entries = append(response.Entries, &SnmpEntry{
			Format: dataUnit.Type,
			Oid:    dataUnit.Name,
			Value:  dataUnit.Value,
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Get is used to retrieve subtree data from a oid using a get request
func Get(p *gosnmp.GoSNMP, oids []string) (*SnmpResponse, error) {
	response := &SnmpResponse{}

	// Keep stats about the duration of the request
	p.OnSent = func(gs *gosnmp.GoSNMP) {
		response.Start = time.Now()
	}
	p.OnRecv = func(gs *gosnmp.GoSNMP) {
		response.End = time.Now()
		response.Duration = response.End.Sub(response.Start)
	}

	// Execute the request
	res, err := p.Get(oids)
	if err != nil {
		return nil, err
	}

	for _, variable := range res.Variables {
		if variable.Type == gosnmp.NoSuchInstance || variable.Type == gosnmp.NoSuchObject {
			continue
		}
		response.Entries = append(response.Entries, &SnmpEntry{
			Format: variable.Type,
			Oid:    variable.Name,
			Value:  variable.Value,
		})
	}

	return response, nil
}

// GetNext is used to retrieve the next value in a subtree
func GetNext(p *gosnmp.GoSNMP, oids []string) (*SnmpResponse, error) {
	response := &SnmpResponse{}

	// Keep stats about the duration of the request
	p.OnSent = func(gs *gosnmp.GoSNMP) {
		response.Start = time.Now()
	}
	p.OnRecv = func(gs *gosnmp.GoSNMP) {
		response.End = time.Now()
		response.Duration = response.End.Sub(response.Start)
	}

	// Execute the request
	res, err := p.GetNext(oids)
	if err != nil {
		return nil, err
	}

	for _, variable := range res.Variables {
		if variable.Type == gosnmp.NoSuchInstance || variable.Type == gosnmp.NoSuchObject || variable.Type == gosnmp.EndOfContents {
			continue
		}
		response.Entries = append(response.Entries, &SnmpEntry{
			Format: variable.Type,
			Oid:    variable.Name,
			Value:  variable.Value,
		})
	}

	return response, nil
}
