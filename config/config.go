package config

type Config struct {
	DB     DBConfig
	Server ServerConfig
	OpenAI OpenAIConfig
}

func LoadConfig() *Config {
	return &Config{
		DB:     LoadDBConfig(),
		Server: LoadServerConfig(),
		OpenAI: LoadOpenAIConfig(),
	}
}
