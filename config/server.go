package config

import "incident-tracker/utils"

type ServerConfig struct {
	Port string
}

func LoadServerConfig() ServerConfig {
	port := utils.GetEnv("APP_PORT", "8088")
	return ServerConfig{
		Port: port,
	}
}
