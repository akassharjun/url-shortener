package model

type Error struct {
	StatusCode int    `json:"status_code"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

type ErrorCode int

const (
	ExternalServiceError ErrorCode = 100
	DatabaseError ErrorCode = 200
)
