package domain

import "time"

type Tweet struct {
	ID        int64
	Author    int64
	Content   string
	CreatedAt time.Time
}
