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
	mappings := make([]*zbxmap.Mapping, len(entries))
	var re *regexp.Regexp

	if opts.TriggerPattern != "" {
		re = regexp.MustCompile(`#INTERFACE`)
	}

	for i := range entries {
		localPattern, remotePattern := buildTriggerPattern(re, `#INTERFACE`, opts, entries[i].LocalPort, entries[i].Port)

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
func buildTriggerPattern(re *regexp.Regexp, reString string, opts *mappingOptions, localPort string, remotePort string) (string, string) {
	var local, remote string

	if opts.TriggerPattern != "" {
		local = replaceInterface(re, reString, opts.TriggerPattern, localPort)
		remote = replaceInterface(re, reString, opts.TriggerPattern, remotePort)
	} else {
		local = ""
		remote = ""
	}

	return local, remote
}

// replaceInterface is used to replace the value of 'reString' in 'p' by 'i' if 'reString' is present in 'p'.
// Otherwise return 'p' as it is.
func replaceInterface(re *regexp.Regexp, reString string, p string, i string) string {
	match := re.Match([]byte(p))
	if match {
		p = strings.Replace(p, reString, i, -1)
	}

	return p
}
