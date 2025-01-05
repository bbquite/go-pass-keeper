package utils

import (
	"os"
	"path/filepath"
)

func SaveFile(filePath string, fileData []byte) error {
	err := os.WriteFile(filePath, fileData, 0666)
	if err != nil {
		return err
	}
	return nil
}

func GetFileInfo(filePath string) ([]byte, string, int64, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", 0, err
	}
	_, fileName := filepath.Split(filePath)

	fi, err := os.Stat(filePath)
	if err != nil {
		return nil, "", 0, err
	}

	fileSize := fi.Size()

	return fileContent, fileName, fileSize, nil
}
