package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tg "github.com/amarnathcjd/gogram/telegram"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Input API ID
	fmt.Print("Enter your API ID: ")
	apiIDStr, _ := reader.ReadString('\n')
	apiIDStr = strings.TrimSpace(apiIDStr)
	var apiID int32
	fmt.Sscanf(apiIDStr, "%d", &apiID)

	// Input API Hash
	fmt.Print("Enter your API Hash: ")
	apiHash, _ := reader.ReadString('\n')
	apiHash = strings.TrimSpace(apiHash)

	// Input Phone Number
	fmt.Print("Enter your phone number (with country code, e.g., +628123456789): ")
	phoneNumber, _ := reader.ReadString('\n')
	phoneNumber = strings.TrimSpace(phoneNumber)

	// Create client configuration with memory session
	cfg := tg.ClientConfig{
		AppID:         apiID,
		AppHash:       apiHash,
		MemorySession: true,
	}

	// Create client
	client, err := tg.NewClient(cfg)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	// Start client (this will prompt for phone and code)
	fmt.Println("\nStarting authentication process...")
	if err := client.Start(); err != nil {
		fmt.Printf("Error starting client: %v\n", err)
		return
	}

	// Check if logged in
	if !client.IsConnected() {
		fmt.Println("Failed to authenticate")
		return
	}

	// Get string session
	stringSession := client.ExportSession()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Your Gogram String Session:")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(stringSession)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nSave this string session securely!")
	fmt.Println("You can use it to login without phone number in the future.")
}
