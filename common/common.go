package common

import (
	"fmt"
	"strconv"
	"time"
	"zlp-demo-golang/config"

	"github.com/google/uuid"
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

// TransID prefix in format: yyMMdd
func GetTransIDPrefix() string {
	now := time.Now()
	return fmt.Sprintf("%02d%02d%02d", now.Year()%100, int(now.Month()), now.Day())
}

// Generate Apptransid in format: yyMMdd_appid_uuidv1
func GenTransID() string {
	return fmt.Sprintf("%v_%v_%v", GetTransIDPrefix(), config.Get("appid"), uuid.New().String())
}
