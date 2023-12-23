package Helpers

import (
	"math/rand"
	"strings"
)

func ContainsInAnyString(searchString string, params ...string) bool {
	for i := range params {
		if strings.Contains(strings.ToLower(params[i]), strings.ToLower(searchString)) {
			return true
		}
	}
	return false
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetRandomString32Lenght() string {
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
