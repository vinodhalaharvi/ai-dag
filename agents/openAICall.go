package agents

import (
	"ai-dag/config"
	"ai-dag/llm"
	"context"
	"fmt"
	"html/template"
	"os"
	"strings"
)

type OpenAICall struct{}

func NewOpenAICall() *OpenAICall {
	return &OpenAICall{}
}

func (o *OpenAICall) Do(
	dagConfig *config.DagConfig,
	agentId string,
	resultCh map[string]chan string,
	childrenResults map[string]string,
) {
	// Create a new GPTChat instance
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		fmt.Println("OPENAI_API_KEY not set")
		return
	}

	t := dagConfig.Agents[agentId]

	// Create a new llm configuration from the agents configuration
	chatConfig := config.ChatConfig{
		Model: t.Model,
	}

	// Add each message from the agents configuration to the llm configuration

	for _, message := range t.Messages {
		parse, err := template.New("content").Parse(message.Content)
		if err != nil {
			fmt.Printf("Failed to parse the message content: %s\n", err)
			return
		}
		strBuilder := &strings.Builder{}

		// Initialize an empty map to hold the data
		data := make(map[string]string)

		// Iterate over the childrenResults map
		for key, value := range childrenResults {
			// Add each key-value pair to the data map
			data[key] = value
		}

		err = parse.Execute(strBuilder, data)
		elems := config.Message{
			Role:    message.Role,
			Content: strBuilder.String(),
		}

		chatConfig.Messages = append(chatConfig.Messages, elems)
		chatConfig.RequestURL = t.URL
		chatConfig.RequestMethod = t.Method
	}

	// Create a new GPTChat instance with the llm configuration
	gptChat := llm.NewGPTChat(key, &config.Config{
		Chat: chatConfig,
	})

	// Execute the llm
	background := context.Background()
	response, err := gptChat.Execute(background, nil)
	if err != nil {
		fmt.Printf("Failed to make the OpenAI API call: %s\n", err)
		return
	}

	// Log the response
	fmt.Printf("OpenAI API response: %s\n", response.(string))

	// Signal this agent's completion
	resultCh[agentId] <- response.(string)
	close(resultCh[agentId])
}
