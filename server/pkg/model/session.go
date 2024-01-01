package model

import "time"

type Session struct {
	Username  string
	Expiry    time.Time
	SessionID string
}
