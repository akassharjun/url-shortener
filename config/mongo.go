package config

type MongoConfig struct {
	Uri         string   `json:"uri"`
	Collections []string `json:"collections"`
}
