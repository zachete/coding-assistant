package config

import "os"

const DefaultModel = "gemini/gemini-3.1-flash-lite"

type Config struct {
	ApiBaseUrl string
	ApiKey     string
	AgentModel string
}

func NewConfig() Config {
	apiKey := os.Getenv("API_KEY")
	apiBaseURL := os.Getenv("API_BASE_URL")
	agentModel := os.Getenv("AGENT_MODEL")

	if apiKey == "" {
		panic("API_KEY environment variable is not set")
	}
	if apiBaseURL == "" {
		apiBaseURL = "https://api.openai.com/v1"
	}
	if agentModel == "" {
		agentModel = DefaultModel
	}

	return Config{
		ApiBaseUrl: apiBaseURL,
		ApiKey:     apiKey,
		AgentModel: agentModel,
	}
}
