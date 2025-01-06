package importkey

import (
	"bufio"
	"fmt"
	"os"
	"pgpcli/internal/listkeys"
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

    fmt.Println()
    fmt.Println("Already stored public key names:")
    currentKeys, err := listkeys.GetPubkeys()
    if err != nil {
        return err
    }
    for _, v := range currentKeys {
        fmt.Println(v)
    }
    fmt.Println()

    fmt.Println("Enter new key name (public key):")
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
