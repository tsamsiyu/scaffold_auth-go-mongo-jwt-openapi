package utils

import "math/rand"

func RandomIntBetween(from int, to int) int {
	return rand.Intn(to-from) + from
}
