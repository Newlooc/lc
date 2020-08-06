package apis

import (
	"time"
)

const (
	DateFormat = "2006-01-02"
)

type URLConfig struct {
	Start time.Time
	End   time.Time
}
