package utility

func GenSetFromSlice(orig []string) map[string]bool {
	var set = make(map[string]bool)
	for _, item := range orig {
		set[item] = true
	}
	return set
}
