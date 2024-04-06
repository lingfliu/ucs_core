package utils

import (
	"time"
)

/**
 * Get current time in nanoseconds
 */
func CurrentTime() int64 {
	t := time.Now()
	return t.UnixNano() / 1000000
}
