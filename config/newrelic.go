package config

type NewRelicConfig struct {
	Licence string `json:"uri"`
	AppName string `json:"appName"`
	Enabled bool   `json:"enabled"`
}
