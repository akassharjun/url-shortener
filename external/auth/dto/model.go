package dto

import "time"

type AuthenticatePayload struct {
	UserId string `json:"user_id"`
	Auth   string `json:"auth"`
}

type AuthenticateResponse struct {
	Validity    time.Time `json:"validity"`
	Permissions struct {
		Read  bool `json:"read"`
		Write bool `json:"write"`
		Admin bool `json:"admin"`
	}
}
