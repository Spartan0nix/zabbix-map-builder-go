package app

import (
	"regexp"
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
	mappings := make([]*zbxmap.Mapping, 0)
	re := regexp.MustCompile(`#INTERFACE`)

	for _, entry := range entries {
		localPattern, remotePattern := buildTriggerPattern(re, opts, entry.LocalPort, entry.Port)

		mappings = append(mappings, &zbxmap.Mapping{
			LocalHost:            localHostname,
			LocalInterface:       entry.LocalPort,
			LocalTriggerPattern:  localPattern,
			LocalImage:           opts.LocalImage,
			RemoteHost:           entry.DeviceId,
			RemoteInterface:      entry.Port,
			RemoteTriggerPattern: remotePattern,
			RemoteImage:          opts.RemoteImage,
		})
	}

	return mappings
}

// buildTriggerPattern is used to build trigger patterns for a mapping.
// Returns local, remote patterns
func buildTriggerPattern(re *regexp.Regexp, opts *mappingOptions, localPort string, remotePort string) (string, string) {
	var local string
	var remote string

	if opts.TriggerPattern != "" {
		local = replaceInterface(re, opts.TriggerPattern, localPort)
		remote = replaceInterface(re, opts.TriggerPattern, remotePort)
	} else {
		local = ""
		remote = ""
	}

	return local, remote
}

// replaceInterface is used to replace the value of 're' in 'p' by 'i' if 're' is present in 'p'.
// Otherwise return 'p' as it is.
func replaceInterface(re *regexp.Regexp, p string, i string) string {
	match := re.FindStringSubmatch(p)

	if len(match) == 0 {
		return p
	}

	return strings.Replace(p, re.String(), i, -1)
}
