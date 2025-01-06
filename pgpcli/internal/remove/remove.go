package remove

import (
	"bufio"
	"fmt"
	"os"
	"pgpcli/internal/listkeys"
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

    pubkeys, err := listkeys.GetPubkeys()
    if err != nil {
        return err
    }
    privkeys, err := listkeys.GetPrivkeys()
    if err != nil {
        return err
    }

    switch u {
    case "1":
        fmt.Println()
        fmt.Println("Stored public keys:")
        for _, v := range pubkeys {
            fmt.Println(v)
        }

        v, err := prompt()
        if err != nil { return err }
        removePubKey(v)
    case "2":
        fmt.Println()
        fmt.Println("Stored private keys:")
        for _, v := range privkeys {
            fmt.Println(v)
        }

        v, err := prompt()
        if err != nil { return err }
        removePrivKey(v)
    case "3": // replace with stored keypairs
        fmt.Println()
        fmt.Println("Keypairs:")
        for _, v := range privkeys {
            for _, w := range pubkeys {
                if w == v {
                    fmt.Println(v)
                }
            }
        }
        fmt.Println()

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
