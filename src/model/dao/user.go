package dao

import (
	"time"
)

type User struct {
	Id        string    `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
