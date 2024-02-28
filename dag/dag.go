package dag

import (
	"ai-dag/agents"
	"ai-dag/config"
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

type DAG struct {
	Lock   sync.Mutex
	Config *config.DagConfig
}

func NewDAG(config *config.DagConfig) *DAG {
	return &DAG{
		Config: config,
		Lock:   sync.Mutex{},
	}
}

func LoadDAGFromYAML(yamlFile string) (*config.DagConfig, error) {
	var cfg config.DagConfig
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (d *DAG) Execute() {
	// Determine execution order
	executionOrder, err := d.topologicalSort()
	if err != nil {
		fmt.Println("Failed to sort agents:", err)
		return
	}

	// Initialize completionCh for all agents
	resultCh := make(map[string]chan string)
	for agentID := range d.Config.Agents {
		resultCh[agentID] = make(chan string)
	}

	for _, agentID := range executionOrder {
		// execute agents in reverse topological order
		data := AgentData{
			AgentId:  agentID,
			ResultCh: resultCh,
		}

		ctx := context.Background()
		go d.executeAgent(ctx, data)
	}

	// Wait for the root agent(s) to complete
	// Assuming the last agent in executionOrder is one of the roots
	<-resultCh[executionOrder[len(executionOrder)-1]]
}

type AgentData struct {
	AgentId  string
	ResultCh map[string]chan string
}

func (d *DAG) executeAgent(ctx context.Context, data AgentData) {
	// Assuming Agents is a map, and you need to access it safely
	d.Lock.Lock()
	children := d.Config.Agents[data.AgentId].Children
	d.Lock.Unlock()

	// Now wait for child agents without holding the lock

	resultCh := data.ResultCh
	childrenResults := make(map[string]string, len(children))
	for _, childID := range children {
		childrenResults[childID] = <-resultCh[childID]
	}

	agentId := data.AgentId
	switch agentId {
	case "fetchCryptoMentions":
		fetchCrypto := agents.NewFetchCryptoMentions()
		fetchCrypto.Do(d.Config, agentId, resultCh, childrenResults)
	case "analyzeCryptoSentiment":
		analyzeSentiment := agents.NewAnalyzeCryptoSentiment()
		analyzeSentiment.Do(d.Config, agentId, resultCh, childrenResults)
	case "nearBySearch":
		agentConfig := d.Config.Agents[agentId]
		request := agents.NewNearBySearchRequest(
			agentConfig.Payload.Location,
			agentConfig.Payload.Radius,
			agentConfig.Payload.Type,
			agentConfig.Payload.Key,
		)
		nearBySearch := agents.NewNearBySearch(request)
		nearBySearch.Do(
			d.Config,
			agentId,
			resultCh,
			childrenResults,
		)
	case "weatherForecast":
		weatherForecast := agents.NewWeatherForecast()
		weatherForecast.Do(d.Config, agentId, resultCh, childrenResults)
	case "openAICall":
		openAICall := agents.NewOpenAICall()
		openAICall.Do(d.Config, agentId, resultCh, childrenResults)
	}
}

func (d *DAG) topologicalSort() ([]string, error) {
	visited := make(map[string]bool)
	result := make([]string, 0)

	var visit func(string) error
	visit = func(nodeID string) error {
		if visited[nodeID] { // Skip already visited nodes
			return nil
		}
		visited[nodeID] = true // Mark node as visited

		// Visit all children
		for _, childID := range d.Config.Agents[nodeID].Children {
			if err := visit(childID); err != nil {
				return err // Handle error as needed
			}
		}

		// Add this node to the result list (post-order)
		result = append(result, nodeID)
		return nil
	}

	// Perform DFS from each node
	for nodeID := range d.Config.Agents {
		if !visited[nodeID] {
			if err := visit(nodeID); err != nil {
				return nil, err // Handle error as needed
			}
		}
	}

	// No need to reverse the result; it's already in topological order
	return result, nil
}
