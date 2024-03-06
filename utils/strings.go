package utils

import "sort"

func IsStringInSlice(s string, slice []string) bool {
	sort.Strings(slice)
	index := sort.SearchStrings(slice, s)
	return index < len(slice) && slice[index] == s
}
