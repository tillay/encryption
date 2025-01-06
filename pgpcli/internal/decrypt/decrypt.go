package decrypt

import (
	"bufio"
	"fmt"
	"os"
	"pgpcli/internal/keyutils"
	"pgpcli/internal/listkeys"
	"pgpcli/lib/clipboard"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"gitlab.com/david_mbuvi/go_asterisks"
)

func Decrypt() error {
    clipText, err := clipboard.Read()
    if err != nil {
        return err
    }

    fmt.Println()
    currentKeys, err := listkeys.GetPrivkeys()
    if err != nil {
        return err
    }
    fmt.Println("Available decrypting keys (private keys):")
    for _, v := range currentKeys {
        fmt.Println(v)
    }
    fmt.Println()

    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println("Enter your private key name:")
    scanner.Scan()
    err = scanner.Err()
    if err != nil {
        return err
    }
    keyname := scanner.Text()

    fmt.Println("Enter your private key passphrase:")
    passphraseBytes, err := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
    if err != nil {
        return err
    }
    passphrase := string(passphraseBytes)

    privKey, err := keyutils.GetPrivKey(passphrase, keyname)
    if err != nil {
        return err
    }

    pgp := crypto.PGP()
    decHandle, err := pgp.Decryption().DecryptionKey(&privKey).New()
    if err != nil {
        return err
    }
    decrypted, err := decHandle.Decrypt([]byte(clipText), crypto.Armor)
    if err != nil {
        return err
    }
    myMessage := string(decrypted.Bytes())

    err = clipboard.Write(myMessage)
    if err != nil {
        return err
    }
    fmt.Println("Decrypted message copied to clipboard!")

    return nil
}
