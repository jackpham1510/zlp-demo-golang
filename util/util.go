package util

import (
	"strconv"
	"time"
)

type Timestamp int64

// Get current time in milliseconds
func GetTimestamp() Timestamp {
	now := time.Now()
	return Timestamp(now.UnixNano() / int64(time.Millisecond))
}

// Convert timestamp to int64
func (t Timestamp) Int() int64 {
	return int64(t)
}

// Convert timestamp to string
func (t Timestamp) String() string {
	return strconv.FormatInt(t.Int(), 10)
}
