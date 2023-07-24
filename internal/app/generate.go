package app

import (
	"strings"

	zbxmap "github.com/Spartan0nix/zabbix-map-builder-go/internal/map"
	"github.com/Spartan0nix/zabbix-map-builder-go/internal/snmp"
)

type mappingOptions struct {
	TriggerPattern string
	LocalImage     string
	RemoteImage    string
}

// generateMapping is used to generate a list of hosts mapping.
func generateMapping(entries []*snmp.SnmpCdpEntry, localHostname string, opts *mappingOptions) []*zbxmap.Mapping {
	mappings := make([]*zbxmap.Mapping, len(entries))

	for i := range entries {
		localPattern, remotePattern := buildTriggerPattern(`#INTERFACE`, opts, entries[i].LocalPort, entries[i].Port)

		mappings[i] = &zbxmap.Mapping{
			LocalHost:            localHostname,
			LocalInterface:       entries[i].LocalPort,
			LocalTriggerPattern:  localPattern,
			LocalImage:           opts.LocalImage,
			RemoteHost:           entries[i].DeviceId,
			RemoteInterface:      entries[i].Port,
			RemoteTriggerPattern: remotePattern,
			RemoteImage:          opts.RemoteImage,
		}

		i++
	}

	return mappings
}

// buildTriggerPattern is used to build trigger patterns for a mapping.
// Returns local, remote patterns
func buildTriggerPattern(pattern string, opts *mappingOptions, localPort string, remotePort string) (string, string) {
	var local, remote string

	if opts.TriggerPattern != "" {
		local = replaceInterface(pattern, opts.TriggerPattern, localPort)
		remote = replaceInterface(pattern, opts.TriggerPattern, remotePort)
	} else {
		local = ""
		remote = ""
	}

	return local, remote
}

// replaceInterface is used to replace the value of 'pattern' in 's' by 'i' if 'pattern' is present in 's'.
// Otherwise return 's' as it is.
func replaceInterface(pattern string, s string, i string) string {
	match := strings.Contains(s, pattern)
	if match {
		s = strings.Replace(s, pattern, i, -1)
	}

	return s
}
