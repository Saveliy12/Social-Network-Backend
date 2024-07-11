package models

import "time"

type Post struct {
	Id            int       `json:"id" db:"id"`
	Content       string    `json:"content" db:"content"`
	Author        string    `json:"author" db:"author"`
	Tags          []string  `json:"tags" db:"tags"`
	CreatedAt     time.Time `json:"createdAt" db:"createdAt"`
	LikesCount    int       `json:"likesCount" db:"likesCount"`
	DislikesCount int       `json:"dislikesCount " db:"dislikesCount"`
}
