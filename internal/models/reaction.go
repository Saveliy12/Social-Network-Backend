package models

import "time"

type Reaction struct {
	Id           int       `json:"id" db:"id"`
	UserId       int       `json:"userId" db:"userId"`
	PostId       int       `json:"postId" db:"postId"`
	ReactionType int       `json:"reactionType" db:"reactionType"`
	CreatedAt    time.Time `json:"createdAt" db:"createdAt"`
}
