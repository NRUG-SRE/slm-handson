package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt returns a random integer between min and max (inclusive)
func RandomInt(min, max int) int {
	if max <= min {
		return min
	}
	return min + rand.Intn(max-min+1)
}

// RandomFloat returns a random float64 between 0.0 and 1.0
func RandomFloat() float64 {
	return rand.Float64()
}

// RandomBool returns a random boolean value
func RandomBool() bool {
	return rand.Intn(2) == 1
}

// RandomSleep sleeps for a random duration between min and max milliseconds
func RandomSleep(minMs, maxMs int) {
	if maxMs <= minMs {
		time.Sleep(time.Duration(minMs) * time.Millisecond)
		return
	}
	
	sleepTime := RandomInt(minMs, maxMs)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
}