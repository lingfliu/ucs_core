package utils

import "math/rand"

func RandInt64(min int64, max int64) int64 {
	return min + rand.Int63()%(max-min)
}

func RandInt32(min int32, max int32) int32 {
	return min + rand.Int31()%(max-min)
}

func RandFloat32(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}
