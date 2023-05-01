package app

import zbxMap "github.com/Spartan0nix/zabbix-map-builder-go/internal/map"

// getUniqueHosts is used to remove dublicate hosts from a Mapping slice.
func getUniqueHosts(mappings []*zbxMap.Mapping) map[string]string {
	out := make(map[string]string, 0)

	for _, m := range mappings {
		_, exist := out[m.LocalHost]
		if !exist {
			out[m.LocalHost] = ""
		}

		_, exist = out[m.RemoteHost]
		if !exist {
			out[m.RemoteHost] = ""
		}
	}

	return out
}
