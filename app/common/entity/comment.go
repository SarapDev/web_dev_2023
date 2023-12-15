package entity

import "time"

type Comment struct {
	ID        int
	PostID    int `db:"post_id"`
	Author    string
	Content   string
	CreatedAt time.Time `db:"created_at"`
}
