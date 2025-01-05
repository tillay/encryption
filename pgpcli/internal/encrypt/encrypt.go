package encrypt

import (
	"bufio"
	"fmt"
	"os"
	"pgpcli/internal/keyutils"
	"pgpcli/lib/clipboard"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func Encrypt() error {
    fmt.Println("Encrypted message:")
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

    for {
        fmt.Println("Name of next recipient, or type HALT:")
        scan := bufio.NewScanner(os.Stdin)
        scan.Scan()
        err := scan.Err()
        if err != nil {
            return err
        }
        next := scan.Text()
        if next == "HALT" {
            break
        }
        nextKey, err := keyutils.GetPubKeyOfUser(next)
        if err != nil {
            return err
        }
        keyring.AddKey(&nextKey)
    }
    encHandle, err := pgp.Encryption().Recipients(keyring).New()

    pgpMessage, err := encHandle.Encrypt([]byte(v))
    armored, err := pgpMessage.Armor()

    err = clipboard.Write(armored)
    if err != nil {
        return err
    }
    fmt.Println("Encrypted message copied to clipboard!")

    return nil
}
