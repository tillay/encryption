package export

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func Export(filename string) error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }

    src, err := os.Open(homeDir + "/wpgp/MAINKEY.pub")
    if err != nil {
        return err
    }
    defer src.Close()

    // Create destination file
    dst, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer dst.Close()

    // Copy source to destination
    _, err = io.Copy(dst, src)
    if err != nil {
        return err
    }

    return nil
}

func HandleExport() error {
	var filepath string
	if len(os.Args) >= 3 {
		filepath = os.Args[2]
	} else {
		fmt.Print("Enter the file path to export the key: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			filepath = scanner.Text()
		}
	}

	if err := Export(filepath); err != nil {
	    return err
    }
	fmt.Printf("Key successfully exported to %s\n", filepath)

    return nil
}
