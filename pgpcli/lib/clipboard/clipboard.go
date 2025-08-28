package clipboard

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func Write(text string) (error) {
    sessionType := os.Getenv("XDG_SESSION_TYPE")
    switch sessionType {
    case "x11":
        return x11Write(text)
    case "wayland":
        return waylandWrite(text)
    default:
        return errors.New("Compositor not supported")
    }
}

func x11Write(text string) error {
    cmd := exec.Command("xclip", "-selection", "clipboard")
    cmd.Stdin = strings.NewReader(text)
    err := cmd.Run()
    return err
}

func waylandWrite(text string) error {
    cmd := exec.Command("wl-copy")
    cmd.Stdin = strings.NewReader(text)

    err := cmd.Run()
    return err
}

func Read() (string, error) {
    sessionType := os.Getenv("XDG_SESSION_TYPE")
    switch sessionType {
    case "x11":
        return x11Read()
    case "wayland":
        return waylandRead()
    default:
        return "", errors.New("Compositor not supported")
    }
}

func x11Read() (string, error) {
    cmd := exec.Command("xclip", "-selection", "clipboard", "-o")
    output, err := cmd.Output()
	if err != nil {
		return "", err
	}
    returnString := string(output)

    return returnString, err
}

func waylandRead() (string, error) {
    cmd := exec.Command("wl-paste")
    output, err := cmd.Output()
	if err != nil {
		return "", err
	}
    returnString := string(output)

    return returnString, err
}
