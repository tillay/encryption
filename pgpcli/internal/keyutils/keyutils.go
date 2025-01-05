package keyutils

import (
	"os"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"golang.design/x/clipboard"
)

func GetMainPrivKey(passphrase string) (crypto.Key, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return crypto.Key{}, err
    }

    lockedKeyFile, err := os.Open(homeDir + "/wpgp/MAINKEY")
    if err != nil {
        return crypto.Key{}, err
    }
    lockedKeyBytes, err := os.ReadFile(lockedKeyFile.Name())
    if err != nil {
        return crypto.Key{}, err
    }
    lockedKeyText := string(lockedKeyBytes)

    lockedKey, err := crypto.NewPrivateKeyFromArmored(lockedKeyText, []byte(passphrase))
    if err != nil {
        return crypto.Key{}, err
    }

    return *lockedKey, err
}

func GetPubKeyOfUser(user string) (crypto.Key, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return crypto.Key{}, err
    }

    pubKeyFile, err := os.Open(homeDir + "/wpgp/" + user + ".pub")
    if err != nil {
        return crypto.Key{}, err
    }
    pubKeyBytes, err := os.ReadFile(pubKeyFile.Name())
    if err != nil {
        return crypto.Key{}, err
    }
    pubKeyText := string(pubKeyBytes)

    pubKey, err := crypto.NewKeyFromArmored(pubKeyText)

    return *pubKey, nil
}

func CheckClipboardForKey() (string, error) {
    err := clipboard.Init()
    if err != nil {
        return "", err
    }
    clipboardBytes := clipboard.Read(clipboard.FmtText)
    clipText := string(clipboardBytes)

    _, err = crypto.NewKeyFromArmored(clipText)
    if err == nil {
        return "import", nil
    }

    if crypto.IsPGPMessage(clipText) {
        return "decrypt", nil
    }

    return "", nil
}
