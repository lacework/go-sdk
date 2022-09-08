package array

import "sort"

func Sort2D(slice [][]string) {
	sort.Slice(slice[:], func(i, j int) bool {
		for x := range slice[i] {
			if slice[i][x] == slice[j][x] {
				continue
			}
			return slice[i][x] < slice[j][x]
		}
		return false
	})
}
