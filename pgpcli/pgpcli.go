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
    if len(os.Args) < 2 {
        helpMessage()
        log.Fatal("No argument provided")
    }
    action := os.Args[1]

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
        fmt.Println(`1. Create key
2. Import key from clipboard
3. Export key to file
4. Encrypt message
5. Decrypt a message
6. List available keys`)
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
        fmt.Println(`./pgpcli create            creates a key
./pgpcli import            imports a key from clipboard
./pgpcli export <filename> exports key to a file
./pgpcli encrypt           encrypt a message
./pgpcli decrypt           decrypts a message from clipboard
./pgpcli list-keys         lists all available pubkeys
./pgpcli help              prints this message`)
}
