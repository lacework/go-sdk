package unique

func StringSlice(slice []string) []string {
	_map := make(map[string]string)

	for _, f := range slice {
		_map[f] = ""
	}

	set := make([]string, len(_map))

	i := 0
	for key := range _map {
		set[i] = key
		i++
	}

	return set
}
