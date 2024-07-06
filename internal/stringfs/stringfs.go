package stringfs

import (
	"fmt"
	"io/fs"
	"os"
)

func SafeRemoveFile(path string, backupSuffix string) error {
	RemoveFile(path + backupSuffix)
	err := RemoveFile(path)
	if err != nil {
		return err
	}

	return nil
}

func SafeWriteFile(path string, backupSuffix string, content string, mode fs.FileMode) error {
	rawContent := []byte(content)

	err := os.WriteFile(
		path,
		rawContent,
		mode,
	)

	if err != nil {
		return err
	}

	err = os.WriteFile(
		path+backupSuffix,
		rawContent,
		mode,
	)

	if err != nil {
		return err
	}

	return nil
}

func SafeReadFile(path string, backupSuffix string) (string, error) {
	rawData, err := os.ReadFile(path)
	if err != nil {
		rawData, err := os.ReadFile(path + backupSuffix)
		if err != nil {
			return "", err
		}
		return string(rawData), nil
	}

	return string(rawData), nil
}

func IsSafeFile(path string, backupSuffix string) (bool, bool) {
	exists, isDir := IsDir(path)
	existsBackup, isDirBackup := IsDir(path + backupSuffix)

	return exists && !isDir,
		existsBackup && !isDirBackup
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
