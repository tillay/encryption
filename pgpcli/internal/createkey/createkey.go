package createkey

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"pgpcli/internal/listkeys"

	"github.com/ProtonMail/gopenpgp/v3/constants"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ProtonMail/gopenpgp/v3/profile"

    "gitlab.com/david_mbuvi/go_asterisks"
)

func CreateKey() error {
    pgp := crypto.PGP()

    fmt.Println()
    currentKeys, err := listkeys.GetPrivkeys()
    if err != nil {
        return err
    }
    fmt.Println("Already stored key names (can be overwritten):")
    for _, v := range currentKeys {
        fmt.Println(v)
    }
    fmt.Println()

    fmt.Println("Enter the new key's passphrase:")
    scanner := bufio.NewScanner(os.Stdin)
    passphrase, err := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
    if err != nil {
        return err
    }
    passphraseBytes := passphrase

    fmt.Println("Enter new key name:")
    scanner.Scan()
    err = scanner.Err()
    if err != nil {
        return err
    }
    keyname := scanner.Text()

    keygenhandle := crypto.PGPWithProfile(profile.RFC9580()).KeyGeneration().AddUserId("createdwithwpgp", "nowhere@goesnowhere.com").New()
    privKey, err := keygenhandle.GenerateKeyWithSecurity(constants.HighSecurity)
    if err != nil {
        return err
    }
    lockedKey, err := pgp.LockKey(privKey, passphraseBytes)
    if err != nil {
        return err
    }

    homeDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }

    _, err = os.Stat(homeDir + "/wpgp")
    if errors.Is(err, os.ErrNotExist) {
        err = os.Mkdir(homeDir + "/wpgp/", 0755)
        if err != nil {
            return err
        }
    }
    err = os.Chmod(homeDir + "/wpgp", 0755)
    if err != nil {
        return err
    }

    _, err = os.Stat(homeDir + "/wpgp/" + keyname)
    if err == nil {
        os.Remove(homeDir + "/wpgp/" + keyname)
    }
    keyFile, err := os.Create(homeDir + "/wpgp/" + keyname)
    if err != nil {
        return err
    }
    defer keyFile.Close()
    keyString, err := lockedKey.Armor()
    if err != nil {
        return err
    }
    _, err = keyFile.WriteString(keyString)
    if err != nil {
        return err
    }

    _, err = os.Stat(homeDir + "/wpgp/" + keyname + ".pub")
    if err == nil {
        os.Remove(homeDir + "/wpgp/" + keyname + ".pub")
    }
    pubKeyFile, err := os.Create(homeDir + "/wpgp/" + keyname + ".pub")
    if err != nil {
        return err
    }
    defer pubKeyFile.Close()
    pubKeyString, err := privKey.GetArmoredPublicKey()
    if err != nil {
        return err
    }
    _, err = pubKeyFile.WriteString(pubKeyString)
    if err != nil {
        return err
    }

    fmt.Println("New key created as " + pubKeyFile.Name())

    return nil
}
