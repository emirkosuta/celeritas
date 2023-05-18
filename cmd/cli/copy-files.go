package main

import (
	"embed"
	"errors"
	"log"
	"os"
	"path/filepath"
)

//go:embed templates
var templateFS embed.FS

func copyFilefromTemplate(templatePath, targetFile string) error {
	if fileExists(targetFile) {
		return errors.New(targetFile + " already exists!")
	}

	data, err := templateFS.ReadFile(templatePath)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile(data, targetFile)
	if err != nil {
		exitGracefully(err)
	}

	return nil
}

func copyDataToFile(data []byte, to string) error {
	dir, _ := filepath.Split(to)                   // get the directory of the target file
	if err := os.MkdirAll(dir, 0755); err != nil { // create the directory if it doesn't exist
		log.Fatalf("failed to create directory: %v", err)
	}

	err := os.WriteFile(to, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func fileExists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}
	return true
}
