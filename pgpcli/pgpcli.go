package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"pgpcli/internal/createkey"
	"pgpcli/internal/decrypt"
	"pgpcli/internal/encrypt"
	"pgpcli/internal/export"
	"pgpcli/internal/importkey"
	"pgpcli/internal/listkeys"
)

func main() {
	var action string
	if len(os.Args) < 2 {
		fmt.Println(
`1. Create key
2. Import key from clipboard
3. Export key to file
4. Encrypt message
5. Decrypt a message
6. List available keys`)
		fmt.Print("Choose action: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			action = scanner.Text()
		}
	} else {
		action = os.Args[1]
	}
	switch action {
		case "1", "create":
			if err := createkey.CreateKey(); err != nil {
				log.Fatalf("Error creating key: %v\n", err)
			}

		case "2", "import":
			if err := importkey.ImportKey(); err != nil {
				log.Fatalf("Error importing key: %v\n", err)
			}

		case "3", "export":
			handleExport()

		case "4", "encrypt":
			if err := encrypt.Encrypt(); err != nil {
				log.Fatalf("Error encrypting message: %v\n", err)
			}

		case "5", "decrypt":
			if err := decrypt.Decrypt(); err != nil {
				log.Fatalf("Error decrypting message: %v\n", err)
			}

		case "6", "list":
			if err := listkeys.ListKeys(); err != nil {
				log.Fatalf("Error listing keys: %v\n", err)
			}

		default:
			fmt.Println("Invalid action. Please try again.")
			helpMessage()
	}
}

func handleExport() {
	var filepath string
	if len(os.Args) >= 3 {
		filepath = os.Args[2]
	} else {
		fmt.Print("Enter the file path to export the key: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			filepath = scanner.Text()
		}
	}

	if err := export.Export(filepath); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Key successfully exported to %s\n", filepath)
}
func helpMessage() {
	fmt.Println(
		`Usage:
		./pgpcli create            Creates a new key
		./pgpcli import            Imports a key from clipboard
		./pgpcli export <filename> Exports a key to a file
		./pgpcli encrypt           Encrypts a message
		./pgpcli decrypt           Decrypts a message from clipboard
		./pgpcli list              Lists all available public keys
		./pgpcli help              Displays this help message`)
}
