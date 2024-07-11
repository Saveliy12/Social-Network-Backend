package models

import "time"

type Friend struct {
	UserId      int       `json:"userId" db:"userId"`
	FriendId    int       `json:"friendId" db:"friendId"`
	FriendLogin string    `json:"friendLogin" db:"friendLogin"`
	AddedAt     time.Time `json:"addedAt" db:"addedAt"`
}
