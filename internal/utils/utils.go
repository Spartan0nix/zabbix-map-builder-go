package utils

// GetMapKey is used to retrieve all the keys from a map of string.
func GetMapKey(m map[string]string) []string {
	out := make([]string, 0)

	for key := range m {
		out = append(out, key)
	}

	return out
}
