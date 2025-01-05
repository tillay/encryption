package importkey

import (
	"bufio"
	"fmt"
	"os"
	"pgpcli/lib/clipboard"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func ImportKey() error {
    clipText, err := clipboard.Read()
    if err != nil {
        return err
    }

    _, err = crypto.NewKeyFromArmored(clipText)
    if err != nil {
        return err
    }

    fmt.Println("Enter username (remember this, its case sensitive!):")
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    err = scanner.Err()
    if err != nil {
        return err
    }
    user := scanner.Text()

    homeDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }

    _, err = os.Stat(homeDir + "/wpgp/" + user + ".pub")
    if !(os.IsNotExist(err)) {
        os.Remove(homeDir + "/wpgp/" + user + ".pub")
    }
    pubKeyFile, err := os.Create(homeDir + "/wpgp/" + user + ".pub")
    if err != nil {
        return err
    }
    defer pubKeyFile.Close()
    pubKeyString := clipText
    if err != nil {
        return err
    }
    _, err = pubKeyFile.WriteString(pubKeyString)
    if err != nil {
        return err
    }

    err = clipboard.Write("")
    if err != nil {
        return err
    }
    return nil
}
