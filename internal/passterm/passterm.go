package passterm

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"syscall"
	"unsafe"
)

func getTermState() (syscall.Termios, error) {
	var termios syscall.Termios

	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	if errno != 0 {
		return syscall.Termios{}, fmt.Errorf("failed to get terminal settings: %v", errno)
	}

	return termios, nil
}

func setTermState(state syscall.Termios) error {
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&state)), 0, 0, 0)

	if errno != 0 {
		return fmt.Errorf("failed to set terminal settings: %v", errno)
	}

	return nil
}

func disableEcho(originalState syscall.Termios) syscall.Termios {
	newState := originalState
	newState.Lflag &^= syscall.ECHO | syscall.ICANON
	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0
	return newState
}

func ReadLine(reader io.Reader) (string, error) {
	bufReader := bufio.NewReader(reader)
	line, err := bufReader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimRight(line, "\r\n"), nil
}

func SilentReadLine(reader io.Reader) (string, error) {
	originalTermios, err := getTermState()

	if err != nil {
		return "", err
	}

	defer setTermState(originalTermios)
	err = setTermState(disableEcho(originalTermios))

	if err != nil {
		return "", err
	}

	bufReader := bufio.NewReader(reader)
	line, err := bufReader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimRight(line, "\r\n"), nil
}
