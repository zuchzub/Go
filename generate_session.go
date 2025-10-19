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

	// Create client configuration
	cfg := tg.NewClientConfigBuilder(apiID, apiHash).
		WithSession("temp_session.dat").
		Build()

	// Create client
	client, err := tg.NewClient(cfg)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	// Connect to Telegram
	_, err = client.Conn()
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		return
	}

	// Login with phone number
	fmt.Println("\nSending verification code...")
	err = client.Login(phoneNumber)
	if err != nil {
		fmt.Printf("Error during login: %v\n", err)
		return
	}

	// Input verification code
	fmt.Print("Enter the verification code: ")
	code, _ := reader.ReadString('\n')
	code = strings.TrimSpace(code)

	// Verify code
	_, err = client.VerifyCode(phoneNumber, code)
	if err != nil {
		// Check if 2FA is enabled
		if strings.Contains(err.Error(), "password") || strings.Contains(err.Error(), "2FA") {
			fmt.Print("Enter your 2FA password: ")
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)

			_, err = client.VerifyPassword(password)
			if err != nil {
				fmt.Printf("Error with 2FA: %v\n", err)
				return
			}
		} else {
			fmt.Printf("Error verifying code: %v\n", err)
			return
		}
	}

	// Get string session
	stringSession := client.ExportStringSession()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Your Gogram String Session:")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(stringSession)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nSave this string session securely!")
	fmt.Println("You can use it to login without phone number in the future.")
	
	// Clean up temporary session file
	os.Remove("temp_session.dat")
}
