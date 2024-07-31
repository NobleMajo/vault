package userin

import (
	"bufio"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

func PromptNewPassword() (string, error) {
	fmt.Println("Enter your new vault password:")
	newPassword, err := ReadPassword()
	if err != nil {
		return "", err
	}
	fmt.Println("Re-enter your new vault password:")

	newPassword2, err := ReadPassword()
	if err != nil {
		return "", err
	}

	if newPassword != newPassword2 {
		return "", fmt.Errorf("passwords don't match")
	}

	return newPassword, nil
}

func PromptPassword() (string, error) {
	fmt.Println("Enter your vault password:")
	newPassword, err := ReadPassword()
	if err != nil {
		return "", err
	}

	return newPassword, nil
}

func ReadPassword() (string, error) {
	rawData, err := term.ReadPassword(int(syscall.Stdin))

	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

func ReadLine() (string, error) {
	fmt.Print("> ")
	buf := bufio.NewReader(os.Stdin)
	rawData, err := buf.ReadBytes('\n')

	if err != nil {
		return "", err
	}

	return string(rawData), nil
}
