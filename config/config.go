package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// Config represents the top-level configuration structure.
type Config struct {
	Chat ChatConfig `yaml:"chat"`
}

// ChatConfig holds configuration specific to llm interactions.
type ChatConfig struct {
	RequestMethod string    `yaml:"requestMethod"`
	RequestURL    string    `yaml:"requestURL"`
	Model         string    `yaml:"model"`
	Messages      []Message `yaml:"messages"`
}

func (c *Config) FromYAMLFile(configPath string) *Config {
	// Read the YAML configuration file
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	// Unmarshal the YAML into the Config struct
	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		os.Exit(1)
	}

	return config
}

// Message structure as defined in assistant.go remains the same
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DagConfig struct {
	Agents map[string]AgentConfig `yaml:"agents"`
}

type Location struct {
	Lat float64 `json:"lat" yaml:"lat"`
	Lng float64 `json:"lng" yaml:"lng"`
}

type AgentConfig struct {
	ID             string    `yaml:"id"`
	Children       []string  `yaml:"children"`
	Plugin         string    `yaml:"agents"`
	PromptTemplate string    `yaml:"promptTemplate"`
	URL            string    `yaml:"url,omitempty"`
	Method         string    `yaml:"method,omitempty"`
	Model          string    `yaml:"model,omitempty"`
	Messages       []Message `yaml:"messages,omitempty"`
	Payload        struct {
		Key      string   `json:"key" yaml:"key"`
		Location Location `json:"location" yaml:"location"`
		Radius   int      `json:"radius" yaml:"radius"`
		Type     string   `json:"type" yaml:"type"`
	} `json:"payload" yaml:"payload"`
	QueryParameters struct {
		Exclude string  `json:"exclude" yaml:"exclude"`
		Lang    string  `json:"lang" yaml:"lang"`
		Lat     float64 `json:"lat" yaml:"lat"`
		Lon     float64 `json:"lon" yaml:"lon"`
		Units   string  `json:"units" yaml:"units"`
	} `json:"queryParameters" yaml:"queryParameters"`
}
