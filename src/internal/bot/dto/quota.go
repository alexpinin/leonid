package dto

import "time"

type Quota struct {
	ChatID        int64
	LastResetDate time.Time
	Remaining     int
}
