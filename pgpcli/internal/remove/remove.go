package remove

import (
	"bufio"
	"fmt"
	"os"
)

func Remove() error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }

    fmt.Println("Key to remove:")
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    err = scanner.Err()
    if err != nil {
        return err
    }
    v := scanner.Text()

    err = os.Remove(homeDir + "/wpgp/" + v + ".pub")
    if err != nil {
        return err
    }

    return nil
}
