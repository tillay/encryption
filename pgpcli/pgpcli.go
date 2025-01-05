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
	"pgpcli/internal/keyutils"
	"pgpcli/internal/listkeys"
	"pgpcli/internal/remove"
)

func main() {
    action, err := keyutils.CheckClipboardForKey()
    if err != nil {
        log.Fatal(err)
    }
	if len(os.Args) > 1 {
        action = os.Args[1]
    }

    processAction(action)
}

func processAction(action string) {
    log.SetFlags(log.LstdFlags)
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
        err := export.HandleExport()
        if err != nil {
            log.Fatal(err)
        }
    case "list-keys", "6":
        err := listkeys.ListKeys()
        if err != nil {
            log.Fatal(err)
        }
    case "remove", "7":
        err := remove.Remove()
        if err != nil {
            log.Fatal(err)
        }
    case "help", "8":
        helpMessage()
    default:
        fmt.Println(`1. creates a key
2. imports a key from clipboard
3. exports a key to file
4. encrypts a message
5. decrypts a message
6. lists available keys
7. removes a key
8. prints a help message`)
        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        err := scanner.Err()
        if err != nil {
            log.Fatal(err)
        }
        processAction(scanner.Text())
    }
}



func helpMessage() {
	fmt.Println(`Usage:
./pgpcli create                       Creates a new key
./pgpcli import                       Imports a key from clipboard
./pgpcli export <filename>            Exports a key to a file
./pgpcli encrypt                      Encrypts a message
./pgpcli decrypt                      Decrypts a message from clipboard
./pgpcli list                         Lists all available public keys
./pgpcli remove                       Remove a key
./pgpcli help                         Displays this help message`)
}
