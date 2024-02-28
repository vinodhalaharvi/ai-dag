package llm

import (
	"ai-dag/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// ChatCompletionResponse to mirror the response structure specific to llm interactions
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int            `json:"index"`
		Message      config.Message `json:"message"`
		LogProbs     interface{}    `json:"logprobs"`
		FinishReason string         `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

// GPTChat to encapsulate llm interactions
type GPTChat struct {
	APIKey   string
	Messages []config.Message
	Client   *http.Client
	Config   *config.Config
}

func (g *GPTChat) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	result, err := g.execute()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// NewGPTChat creates a new instance of GPTChat with initialized values.
func NewGPTChat(apiKey string, cfg *config.Config) *GPTChat {
	return &GPTChat{
		APIKey:   apiKey,
		Messages: cfg.Chat.Messages,
		Client:   &http.Client{},
		Config:   cfg,
	}
}

// AddMessage adds a new message to the conversation.
func (g *GPTChat) AddMessage(role, content string) {
	g.Messages = append(g.Messages, config.Message{Role: role, Content: content})
}

// Execute sends the conversation to OpenAI's API and returns the llm's response.
func (g *GPTChat) execute() (string, error) {
	requestData := map[string]interface{}{
		"model":    g.Config.Chat.Model,
		"messages": g.Messages,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		g.Config.Chat.RequestMethod,
		g.Config.Chat.RequestURL,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+g.APIKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := g.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response ChatCompletionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from OpenAI")
}
