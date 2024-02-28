package utils

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func ExecuteCommandInBash(commandStr string) []byte {
	// Use build.sh to execute the command string
	out, err := exec.Command("bash", "-c", commandStr).CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return nil
	}

	// Return the output
	return out
}

// ExecuteCommand executes a shell command and prints its output

func ToPrettyJsonFromObject(response interface{}) string {
	// Pretty print the JSON response
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(jsonData)
}
