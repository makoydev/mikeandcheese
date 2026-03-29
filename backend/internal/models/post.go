package models

import "time"

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Excerpt   string    `json:"excerpt"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Category  string    `json:"category"`
	Tags      string    `json:"tags"`
	Featured  bool      `json:"featured"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
