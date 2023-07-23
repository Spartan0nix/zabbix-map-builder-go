package snmp

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Spartan0nix/zabbix-map-builder-go/internal/logging"
	"github.com/gosnmp/gosnmp"
)

type SnmpResponse struct {
	Start    time.Time
	End      time.Time
	Duration time.Duration
	Entries  []*SnmpEntry
}

type SnmpEntry struct {
	Format gosnmp.Asn1BER
	Oid    string
	Value  interface{}
}

type SnmpCdpEntry struct {
	Address        []uint8
	DeviceId       string
	Port           string
	Capabilities   []uint8
	LocalPortIndex string
	LocalPort      string
}

// LogRequestDuration is used to log the time took by the request.
func LogRequestDuration(logger *logging.Logger, d time.Duration) {
	logger.Debug(fmt.Sprintf("Request took : %s", d))
}

// SupportedCdpCacheCapabilities is used to check if the given interface capabilities are those of a Router or Switch device.
func SupportedCdpCacheCapabilities(v []uint8) bool {
	switch fmt.Sprintf("%x", v[len(v)-1]) {
	case "28":
		return true
	case "29":
		return true
	default:
		return false
	}
}

// ParseCdpCache is used to consolidate a given list of individual entries from a CdpCacheTable.
func ParseCdpCache(entries []*SnmpEntry, logger *logging.Logger) []*SnmpCdpEntry {
	// Group 1 : .
	// Group 2 : MIB Node (6 -> cdpCacheDeviceId, 7 -> cdpCacheDevicePort, etc.)
	// Group 3 : index used to retrieve local interface value (IF-MIB)
	re := regexp.MustCompile(`(.)?1.3.6.1.4.1.9.9.23.1.2.1.1.(\d+).(\d+).\d+`)
	cdpEntries := make([]*SnmpCdpEntry, 0)

	// indexExcluded keep track of entry without enough capacilities to be consider a router or a switch
	indexExcluded := make([]int, 0)

	for _, entry := range entries {
		match := re.FindSubmatch([]byte(entry.Oid))

		if len(match) != 4 {
			logger.Warning(fmt.Sprintf("error while parsing oid '%s'", entry.Oid))
			logger.Debug(fmt.Sprintf("Regex used to extract data from each oid : %s", re))
			logger.Debug("The second captured group is used to extract the node name (cdpCacheDeviceId, cdpCacheDevicePort, etc.)")
			logger.Debug("The third captured group is used to extract the interface index")

			continue
		}

		node := string(match[2])
		index := string(match[3])
		// Keep track of the position of the current entry in the cache
		var cachePosition int

		// Check if a struct was already provided for this index or not
		exist := false
		for i, entry := range cdpEntries {
			if entry.LocalPortIndex == index {
				exist = true
				cachePosition = i
			}
		}

		if !exist {
			cdpEntries = append(cdpEntries, &SnmpCdpEntry{
				LocalPortIndex: index,
			})

			cachePosition = len(cdpEntries) - 1
		}

		// Loop over each MIB node index
		switch node {
		// cdpCacheAddress
		case "4":
			cdpEntries[cachePosition].Address = entry.Value.([]uint8)
		// cdpCacheDeviceId
		case "6":
			cdpEntries[cachePosition].DeviceId = string(entry.Value.([]byte))
		// cdpCacheDevicePort
		case "7":
			cdpEntries[cachePosition].Port = string(entry.Value.([]byte))
		// cdpCacheCapabilities
		case "9":
			c := entry.Value.([]uint8)

			// Check for supported capabilities
			if ok := SupportedCdpCacheCapabilities(c); ok {
				// If the device capabilities are supported, add them for further use
				cdpEntries[cachePosition].Capabilities = c
			} else {
				// Otherwise, add the host index to the list of entry to be removed
				indexExcluded = append(indexExcluded, cachePosition)
			}
		default:
			logger.Warning(fmt.Sprintf("unsupported SNMP node 'oid : %s' found in the CDP cache", entry.Oid))
		}
	}

	// Remove devices without the required capabilities
	for _, i := range indexExcluded {
		// Overwrite the ptr with the last ptr of the list
		cdpEntries[i] = cdpEntries[len(cdpEntries)-1]
		// Remove the last ptr from the list
		cdpEntries = cdpEntries[:len(cdpEntries)-1]
	}

	return cdpEntries
}

// GetLocalInterfacesName is used to retrieve the name of each local interface associated with a list of cdp entries.
func GetLocalInterfacesName(p *gosnmp.GoSNMP, entries []*SnmpCdpEntry, logger *logging.Logger) error {
	oids := make([]string, 0)
	// Keep track of the position of each interface index in the list
	indexMapping := make(map[string]int, len(entries))

	// Build the list of oids to query for each interface index
	for i, entry := range entries {
		oids = append(oids, fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.1.%s", entry.LocalPortIndex))
		indexMapping[entry.LocalPortIndex] = i
	}

	// Query the oids
	res, err := Get(p, oids)
	if err != nil {
		return err
	}

	LogRequestDuration(logger, res.Duration)

	// Store data as a map : interface index -> name
	for _, entry := range res.Entries {
		currentIndex := string(entry.Oid[len(entry.Oid)-1])
		entries[indexMapping[currentIndex]].LocalPort = string(entry.Value.([]uint8))
	}

	return nil
}

// GetLocalHostname is used to retrieve the name of the host
func GetHostname(p *gosnmp.GoSNMP, logger *logging.Logger) (string, error) {
	res, err := GetNext(p, []string{
		"1.3.6.1.4.1.9.9.23.1.3.4",
	})
	if err != nil {
		return "", err
	}

	LogRequestDuration(logger, res.Duration)

	if len(res.Entries) != 1 {
		return "", fmt.Errorf("wrong number of variables returned while retrieving local hostname. Expected : 1. Returned : %d", len(res.Entries))
	}

	str := string(res.Entries[0].Value.([]uint8))

	return str, nil
}
