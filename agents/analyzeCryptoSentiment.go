package agents

import (
	"ai-dag/config"
	"fmt"
)

type AnalyzeCryptoSentiment struct {
}

func NewAnalyzeCryptoSentiment() *AnalyzeCryptoSentiment {
	return &AnalyzeCryptoSentiment{}
}

func (a *AnalyzeCryptoSentiment) Do(
	config *config.DagConfig,
	agentId string,
	resultCh map[string]chan string,
	childResults map[string]string,
) {
	fmt.Printf("analyzeCryptoSentiment: %s\n", agentId)
	// TODO: Mock implementation
	// you know what to do
	resultCh[agentId] <- "Sentiment: positive"
	close(resultCh[agentId])
}
