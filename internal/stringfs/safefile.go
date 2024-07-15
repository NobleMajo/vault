package stringfs

import (
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

func SafeReadBothFiles(path string, backupSuffix string) (string, error, string, error) {
	rawData1, err1 := os.ReadFile(path)
	rawData2, err2 := os.ReadFile(path + backupSuffix)

	return string(rawData1), err1, string(rawData2), err2
}

func IsSafeFile(path string, backupSuffix string) (bool, bool) {
	exists, isDir := IsDir(path)
	existsBackup, isDirBackup := IsDir(path + backupSuffix)

	return exists && !isDir,
		existsBackup && !isDirBackup
}
