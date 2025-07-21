package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	ok := scanner.Scan()
	if !ok {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", scanner.Err())
		return
	}

	token := strings.TrimSpace(scanner.Text())

	if token == "" {
		return
	}

	// Split the JWT token into parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		fmt.Fprintf(os.Stderr, "Error: Invalid JWT token format (expected 3 parts, got %d)\n", len(parts))
		return
	}

	// Decode header
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding header: %v\n", err)
		return
	}

	// Decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding payload: %v\n", err)
		return
	}

	// Pretty print header
	var header map[string]interface{}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing header JSON: %v\n", err)
		return
	}

	// Pretty print payload
	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing payload JSON: %v\n", err)
		return
	}

	payloadJSON, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting payload: %v\n", err)
		return
	}

	fmt.Println(string(payloadJSON))
}
