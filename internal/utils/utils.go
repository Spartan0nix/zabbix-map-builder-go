package utils

// GetMapKey is used to retrieve all the keys from a map of string.
func GetMapKey(m map[string]string) []string {
	out := make([]string, len(m))
	i := 0

	for key := range m {
		out[i] = key
		i++
	}

	return out
}
