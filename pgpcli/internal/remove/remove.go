package remove

import (
	"bufio"
	"fmt"
	"os"
)

func Remove() error {
    scanner := bufio.NewScanner(os.Stdin)

    fmt.Println(`1. Remove a pubkey
2. Remove a privkey
3. Remove a keypair
Else. abort`)
    scanner.Scan()
    err := scanner.Err()
    if err != nil {
        return err
    }
    u := scanner.Text()
    switch u {
    case "1":
        v, err := prompt()
        if err != nil { return err }
        removePubKey(v)
    case "2":
        v, err := prompt()
        if err != nil { return err }
        removePrivKey(v)
    case "3":
        v, err := prompt()
        if err != nil { return err }
        removePubKey(v)
        removePrivKey(v)
    }

    return nil
}

func removePubKey(keyname string) error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }

    err = os.Remove(homeDir + "/wpgp/" + keyname + ".pub")
    if err != nil {
        return err
    }

    return nil
}

func removePrivKey(keyname string) error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }

    err = os.Remove(homeDir + "/wpgp/" + keyname)
    if err != nil {
        return err
    }

    return nil
}

func prompt() (string, error) {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println("Key to remove:")
    scanner.Scan()
    err := scanner.Err()
    if err != nil {
        return "", err
    }
    v := scanner.Text()
    return v, nil
}
