package stringfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func ParsePath(path *string) error {
	if path == nil {
		return errors.New("path is nil")
	}

	*path = strings.TrimSpace(*path)

	if strings.HasPrefix(*path, "~") {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return errors.New("cant get users home dir:\n> " + err.Error())
		}

		*path = strings.Replace(*path, "~", userHome, 1)
	}

	if !strings.HasPrefix(*path, "/") {
		wd, err := os.Getwd()
		if err != nil {
			return errors.New("cant get current working dir:\n> " + err.Error())
		}

		*path = wd + "/" + *path
	}

	return nil
}

func RemoveFile(path string) error {
	exists, isDir := IsDir(path)

	if !exists {
		return nil
	}

	if isDir {
		return os.RemoveAll(path)
	}

	return os.Remove(path)
}

func WriteFile(path string, content string, mode fs.FileMode) error {
	err := os.WriteFile(
		path,
		[]byte(content),
		mode,
	)

	if err != nil {
		return err
	}

	return nil
}

func ReadFile(path string) (string, error) {
	rawData, err := os.ReadFile(path)
	if err != nil {
		fmt.Print(err)
	}

	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

func Exists(path string) bool {
	stat, err := os.Stat(path)

	if err != nil || stat == nil {
		return false
	}

	return true
}

func IsFile(path string) (bool, bool) {
	stat, err := os.Stat(path)

	if err != nil || stat == nil {
		return false, false
	}

	return true, !stat.IsDir()
}

func IsDir(path string) (bool, bool) {
	stat, err := os.Stat(path)

	if err != nil || stat == nil {
		return false, false
	}

	return true, stat.IsDir()
}
