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
    action := ""
	if len(os.Args) > 1 {
        	action = os.Args[1]
    }
    processAction(action)
}

func processAction(action string) {
    switch action {
    case "create", "1":
        err := createkey.CreateKey()
        if err != nil {
            log.Fatal(err)
        }
    case "encrypt", "4":
        err := encrypt.Encrypt()
        if err != nil {
            log.Fatal(err)
        }
    case "import", "2":
        err := importkey.ImportKey()
        if err != nil {
            log.Fatal(err)
        }
    case "decrypt", "5":
        err := decrypt.Decrypt()
        if err != nil {
            log.Fatal(err)
        }
    case "export", "3":
        if len(os.Args) < 3 {
            log.Fatal("Put in a filepath silly!")
        }
        err := export.Export(os.Args[2])
        if err != nil {
            log.Fatal(err)
        }
    case "list-keys", "6":
        err := listkeys.ListKeys()
        if err != nil {
            log.Fatal(err)
        }
    case "help":
        helpMessage()
    default:
        fmt.Println(`1. creates a key
2. imports a key from clipboard
3. exports a key to file
4. encrypts a message
5. decrypts a message
6. lists available keys
7. prints a help message`)
        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        err := scanner.Err()
        if err != nil {
            log.Fatal(err)
        }
        processAction(scanner.Text())
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
