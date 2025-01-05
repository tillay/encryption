package listkeys

import (
	"fmt"
	"os"
	"strings"
)

func ListKeys() error {
    pubkeys, err := GetPubkeys()
    if err != nil {
        return err
    }
    privkeys, err := GetPrivkeys()
    if err != nil {
        return err
    }

    fmt.Println("Pubkeys:")
    for _, v := range pubkeys {
        fmt.Println(v)
    }
    fmt.Println("\nPrivkeys:")
    for _, v := range privkeys {
        fmt.Println(v)
    }

    return nil
}

func GetPubkeys() ([]string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return make([]string, 0), err
    }

    entries, err := os.ReadDir(homeDir + "/wpgp/")
    if err != nil {
        return make([]string, 0), err
    }
    pubkeys := make([]string, 0, len(entries))

    for _, e := range entries {
        if (strings.Contains(e.Name(), ".pub")) {
            pubkeys = append(pubkeys, strings.Trim(e.Name(), ".pub"))
        }
    }

    return pubkeys, nil
}

func GetPrivkeys() ([]string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return make([]string, 0), err
    }

    entries, err := os.ReadDir(homeDir + "/wpgp/")
    if err != nil {
        return make([]string, 0), err
    }
    privkeys := make([]string, 0, len(entries))

    for _, e := range entries {
        if (!strings.Contains(e.Name(), ".pub")) {
            privkeys = append(privkeys, strings.Trim(e.Name(), ".pub"))
        }
    }

    return privkeys, nil
}
