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
	dir, file := filepath.Split(path)
	rawContent := []byte(content)

	err := os.WriteFile(
		dir+".tmp_"+file,
		rawContent,
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
