package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
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

	// Create client
	client, err := telegram.NewClient(telegram.ClientConfig{
		AppID:    apiID,
		AppHash:  apiHash,
		LogLevel: telegram.LogInfo,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	// Connect
	if err := client.Connect(); err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		return
	}

	// Send code
	fmt.Println("\nSending verification code...")
	sentCode, err := client.Auth().SendCode(phoneNumber, 5)
	if err != nil {
		fmt.Printf("Error sending code: %v\n", err)
		return
	}

	// Input verification code
	fmt.Print("Enter the verification code: ")
	code, _ := reader.ReadString('\n')
	code = strings.TrimSpace(code)

	// Sign in
	_, err = client.Auth().SignIn(phoneNumber, code, sentCode.PhoneCodeHash)
	if err != nil {
		// Check if 2FA is enabled
		if strings.Contains(err.Error(), "SESSION_PASSWORD_NEEDED") {
			fmt.Print("Enter your 2FA password: ")
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)

			_, err = client.Auth().Password(context.Background(), password)
			if err != nil {
				fmt.Printf("Error with 2FA: %v\n", err)
				return
			}
		} else {
			fmt.Printf("Error signing in: %v\n", err)
			return
		}
	}

	// Get string session
	stringSession, err := client.ExportStringSession()
	if err != nil {
		fmt.Printf("Error exporting session: %v\n", err)
		return
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Your Gogram String Session:")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(stringSession)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nSave this string session securely!")
	fmt.Println("You can use it to login without phone number in the future.")
}
