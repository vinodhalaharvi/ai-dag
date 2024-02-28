package agents

import (
	"ai-dag/config"
	"fmt"
)

type FetchCryptoMentions struct {
}

func NewFetchCryptoMentions() *FetchCryptoMentions {
	return &FetchCryptoMentions{}
}

func (f FetchCryptoMentions) Do(
	config *config.DagConfig,
	agentId string,
	resultCh map[string]chan string,
	childResults map[string]string,
) {
	fmt.Printf("fetchCryptoMentions: %s\n", agentId)
	// TODO: Mock implementation
	// you know what to do
	resultCh[agentId] <- "BTC, ETH, SOL"
	close(resultCh[agentId])
}
