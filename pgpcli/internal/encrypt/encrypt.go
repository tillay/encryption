package encrypt

import (
	"bufio"
	"fmt"
	"os"
	"pgpcli/internal/keyutils"
	"pgpcli/internal/listkeys"
	"pgpcli/lib/clipboard"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func Encrypt() error {
    fmt.Println("Message to encrypt:")
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    err := scanner.Err()
    if err != nil {
        return err
    }
    v := scanner.Text()

    pgp := crypto.PGP()

    keyring, err := crypto.NewKeyRing(nil)
    if err != nil {
        return err
    }

    fmt.Println()
    fmt.Println("Stored recepient options (public keys):")
    currentKeys, err := listkeys.GetPubkeys()
    if err != nil {
        return err
    }
    for _, v := range currentKeys {
        fmt.Println(v)
    }
    fmt.Println()


    for {
        fmt.Println("Name of next recipient, or type q:")
        scan := bufio.NewScanner(os.Stdin)
        scan.Scan()
        err := scan.Err()
        if err != nil {
            return err
        }
        next := scan.Text()
        if next == "q" {
            break
        }
        nextKey, err := keyutils.GetPubKeyOfUser(next)
        if err != nil {
            fmt.Println("That didn't work, try again.")
            continue
        }
        keyring.AddKey(&nextKey)
    }
    encHandle, err := pgp.Encryption().Recipients(keyring).New()
    if err != nil {
        return err
    }
    pgpMessage, err := encHandle.Encrypt([]byte(v))
    if err != nil {
        return err
    }
    armored, err := pgpMessage.Armor()
    if err != nil {
        return err
    }

    err = clipboard.Write(armored)
    if err != nil {
        return err
    }
    fmt.Println("Encrypted message copied to clipboard!")

    return nil
}
