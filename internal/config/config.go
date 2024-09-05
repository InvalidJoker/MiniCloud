package config

type Config struct {
	AuthToken string `json:"auth_token"`
	Port      int    `json:"port"`
	Interface string `json:"interface"`
	DatabaseURL string `json:"database_url"`
}
