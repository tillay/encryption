package export

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"pgpcli/internal/createkey"
	"pgpcli/internal/listkeys"
)

func Export(filename string) error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }

    pubKeys, err := listkeys.GetPubkeys()
    if err != nil {
        return err
    }

    if len(pubKeys) == 0 {
        fmt.Println("No key found, creating new key")
        err = createkey.CreateKey()
        if err != nil {
            return err
        }
    }

    fmt.Println()
    fmt.Println("Current stored public keys are:")
    currentKeys, err := listkeys.GetPubkeys()
    if err != nil {
        return err
    }
    for _, v := range currentKeys {
        fmt.Println(v)
    }
    fmt.Println()

    fmt.Println("Key name:")
    scan := bufio.NewScanner(os.Stdin)
    scan.Scan()
    err = scan.Err()
    if err != nil {
        return err
    }
    keyname := scan.Text()

    src, err := os.Open(homeDir + "/wpgp/" + keyname + ".pub")
    if err != nil {
        return err
    }
    defer src.Close()

    // Create destination file
    dst, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer dst.Close()

    // Copy source to destination
    _, err = io.Copy(dst, src)
    if err != nil {
        fmt.Println("here")
        return err
    }

    return nil
}

func HandleExport() error {
	var filepath string
	if len(os.Args) >= 3 {
		filepath = os.Args[2]
	} else {
		fmt.Print("Enter the file path or name to export the key to: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			filepath = scanner.Text()
		}
	}

	if err := Export(filepath); err != nil {
	    return err
    }
	fmt.Printf("Key successfully exported to %s\n", filepath)

    return nil
}
