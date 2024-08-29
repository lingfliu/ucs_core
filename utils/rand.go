package utils

import "math/rand"

func RandInt64(min int64, max int64) int64 {
	return min + rand.Int63()%(max-min)
}
