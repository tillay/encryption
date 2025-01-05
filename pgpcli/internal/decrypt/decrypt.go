package decrypt

import (
	"bufio"
	"fmt"
	"os"
	"pgpcli/internal/keyutils"
	"pgpcli/lib/clipboard"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func Decrypt() error {
    clipText, err := clipboard.Read()
    if err != nil {
        return err
    }

    fmt.Println("Enter key passphrase:")
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    err = scanner.Err()
    if err != nil {
        return err
    }
    passphrase := scanner.Text()

    fmt.Println("Enter key name:")
    scanner.Scan()
    err = scanner.Err()
    if err != nil {
        return err
    }
    keyname := scanner.Text()

    privKey, err := keyutils.GetPrivKey(passphrase, keyname)
    if err != nil {
        return err
    }

    pgp := crypto.PGP()
    decHandle, err := pgp.Decryption().DecryptionKey(&privKey).New()
    decrypted, err := decHandle.Decrypt([]byte(clipText), crypto.Armor)
    myMessage := string(decrypted.Bytes())

    err = clipboard.Write(myMessage)
    if err != nil {
        return err
    }
    fmt.Println("Encrypted message copied to clipboard!")

    return nil
}
