package utils

import (
	"time"
)

/**
 * Get current time in nanoseconds
 */
func CurrentTime() int64 {
	return time.Now().UnixNano()
}
