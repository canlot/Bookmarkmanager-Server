package Helpers

import "strings"

func ContainsInAnyString(searchString string, params ...string) bool {
	for i := range params {
		if strings.Contains(strings.ToLower(params[i]), strings.ToLower(searchString)) {
			return true
		}
	}
	return false
}
