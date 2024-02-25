package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"sync"
	"time"
)

type Task struct {
	ID       string   `yaml:"id"`
	Children []string `yaml:"children"`
}

type DAG struct {
	Tasks map[string]Task `yaml:"tasks"`
	Lock  sync.Mutex
}

func main() {
	dag, err := LoadDAGFromYAML("tasks.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Determine execution order
	executionOrder, err := topologicalSort(dag)
	if err != nil {
		fmt.Println("Failed to sort tasks:", err)
		return
	}

	// Initialize channels for all tasks
	channels := make(map[string]chan bool)
	for taskID := range dag.Tasks {
		channels[taskID] = make(chan bool)
	}

	// Execute tasks in reverse topological order
	for _, taskID := range executionOrder {
		go executeTask(taskID, dag, channels)
	}

	// Wait for the root task(s) to complete
	// Assuming the last task in executionOrder is one of the roots
	<-channels[executionOrder[len(executionOrder)-1]]
	fmt.Printf("DAG execution completed\n")
}

func LoadDAGFromYAML(yamlFile string) (*DAG, error) {
	var dag DAG
	data, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &dag)
	if err != nil {
		return nil, err
	}
	return &dag, nil
}

func executeTask(taskID string, dag *DAG, channels map[string]chan bool) {
	// Assuming dag.Tasks is a map and you need to access it safely
	dag.Lock.Lock()
	children := dag.Tasks[taskID].Children
	dag.Lock.Unlock()

	// Now wait for child tasks without holding the lock
	for _, childID := range children {
		<-channels[childID]
	}

	// Perform the task's work (simulate with a print statement)
	// Create a ticker that ticks every second.
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Create a timer that fires after 5 seconds.
	timer := time.NewTimer(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			// This block executes every second.
			fmt.Println("Tick at", time.Now())

		case <-timer.C:
			// This block executes after 5 seconds.
			fmt.Println("Timer expired, exiting...")
			// Signal this task's completion
			close(channels[taskID])
			return
		}
	}
}

func topologicalSort(dag *DAG) ([]string, error) {
	visited := make(map[string]bool)
	result := make([]string, 0)

	var visit func(string) error
	visit = func(nodeID string) error {
		if visited[nodeID] { // Skip already visited nodes
			return nil
		}
		visited[nodeID] = true // Mark node as visited

		// Visit all children
		for _, childID := range dag.Tasks[nodeID].Children {
			if err := visit(childID); err != nil {
				return err // Handle error as needed
			}
		}

		// Add this node to the result list (post-order)
		result = append(result, nodeID)
		return nil
	}

	// Perform DFS from each node
	for nodeID := range dag.Tasks {
		if !visited[nodeID] {
			if err := visit(nodeID); err != nil {
				return nil, err // Handle error as needed
			}
		}
	}

	// No need to reverse the result; it's already in topological order
	return result, nil
}
