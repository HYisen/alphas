package utility

import "os"

func GenSetFromSlice(orig []string) map[string]bool {
	var set = make(map[string]bool)
	for _, item := range orig {
		set[item] = true
	}
	return set
}

func Filter(source []string, predictor func(string) bool) []string {
	var ret []string
	for _, one := range source {
		if predictor(one) {
			ret = append(ret, one)
		}
	}
	return ret
}

// AcquireNoFlagArgs use os.Args as source, remove first exe path and any -(hyphen) started parameters.
// Lacking information of flags, it can not handle flags like "-XXX" / "-XXX true".
func AcquireNoFlagArgs() []string {
	return Filter(os.Args[1:], func(s string) bool {
		// As arg can not be empty string, it's redundant to use !strings.HasPrefix(s, "=").
		return s[0] != '-'
	})
}
