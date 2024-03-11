package dao

import "time"

type Link struct {
	Id        string    `json:"link_id"`
	UserId    *string   `json:"user_id"`
	ShortURL  string    `json:"short_url"`
	LongURL   string    `json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
}

type LinkQuery struct {
	Id        *string    `bson:"link_id"`
	UserId    *string    `bson:"user_id"`
	ShortURL  *string    `bson:"short_url"`
	LongURL   *string    `bson:"long_url"`
	CreatedAt *time.Time `bson:"created_at"`
}
