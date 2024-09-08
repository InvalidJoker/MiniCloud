package config

type Config struct {
	AuthToken    string `json:"auth_token"`
	Port         int    `json:"port"`
	DatabaseType string `json:"database_type"`
	DatabaseURL  string `json:"database_url"`
}
