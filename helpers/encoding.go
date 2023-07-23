package helpers

import (
	"math/rand"
)

// Base62
// [0-9] = 10 Characters
// [A-Z] = 26 Characters
// [a-z] = 26 Characters
// Total = 62 Characters

const (
	character_set = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	link_length   = 7
)

func EncodeBase62() string {
	var encodedString string

	// rand.Intn(62) // Generate Random Integer between 0 - 61
	for i := 0; i < link_length; i++ {
		encodedString += string(character_set[rand.Intn(len(character_set))])
	}
	return encodedString
}
