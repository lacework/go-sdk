package array

import (
	"sort"
	"strings"
)

// Sort2D can be used 2d Arrays used for Table Outputs to sort headers
func Sort2D(slice [][]string) {

	for range slice {
		sort.Slice(slice[:], func(i, j int) bool {
			elem := slice[i][0]
			next := slice[j][0]
			switch strings.Compare(elem, next) {
			case -1, 0:
				return true
			case 1:
				return false
			}
			return false
		})
	}
}

func SortStrings(slice []string) []string {
	sort.Strings(slice)

	return slice
}
