package url

import (
	"sync/atomic"
	"time"
)

type URL struct {
	ID          string
	UserID      string
	LongURL     string
	ShortCode   string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	IsActive    bool
	IsDeleted   bool
	IsCustom    bool
	CustomAlias string
	clicks      atomic.Uint64
}

// Clicks returns the number of successful resolves (redirects) recorded for this URL.
func (u *URL) Clicks() uint64 {
	if u == nil {
		return 0
	}
	return u.clicks.Load()
}

func (u *URL) recordClick() {
	u.clicks.Add(1)
}
