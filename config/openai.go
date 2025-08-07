package config

import "incident-tracker/utils"

type OpenAIConfig struct {
	APIKey      string
	APIUrl      string
	Model       string
	Temperature float64
}

func LoadOpenAIConfig() OpenAIConfig {
	return OpenAIConfig{
		APIKey:      utils.GetEnv("OPEN_AI_API_KEY", ""),
		APIUrl:      utils.GetEnv("OPENAI_API_URL", "https://api.openai.com/v1/chat/completions"),
		Model:       utils.GetEnv("OPENAI_MODEL", "gpt-3.5-turbo"),
		Temperature: utils.GetEnvFloat("OPENAI_TEMPERATURE", 0.2),
	}
}
