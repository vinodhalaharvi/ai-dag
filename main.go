package main

import (
	"ai-dag/dag"
)

func main() {
	config, err := dag.LoadDAGFromYAML("graph.yaml")
	if err != nil {
		panic(err)
	}

	// Load dGraph into registry
	dGraph := dag.NewDAG(config)
	dGraph.Execute()
}
