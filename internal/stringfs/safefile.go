package stringfs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

func RemoveTmpSafeFile(path string) {
	dir, file := filepath.Split(path)

	RemoveFile(dir + ".tmp_" + file)
}

func SafeWriteFile(path string, content string, mode fs.FileMode) error {
	return SafeWriteFileBytes(path, []byte(content), mode)
}

func SafeWriteFileBytes(path string, content []byte, mode fs.FileMode) error {
	dir, file := filepath.Split(path)

	err := os.WriteFile(
		dir+".tmp_"+file,
		content,
		mode,
	)

	if err != nil {
		return errors.New("Write file error: " + err.Error())
	}

	err = os.Rename(dir+".tmp_"+file, path)

	if err != nil {
		return errors.New("Rename file error: " + err.Error())
	}

	return nil
}
