package stringfs

import (
	"io/fs"
	"os"
)

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
