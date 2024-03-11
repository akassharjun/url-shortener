package model

type Response struct {
	StatusCode int         `json:"status_code"`
	Code       int         `json:"code"`
	Data       interface{} `json:"data"`
}
