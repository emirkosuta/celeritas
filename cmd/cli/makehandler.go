package main

import (
	"errors"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
)

func makeHandler(arg3 string) error {
	if arg3 == "" {
		return errors.New("you must give the handler a name")
	}

	fileName := cel.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
	if fileExists(fileName) {
		return errors.New(fileName + " already exists!")
	}

	data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
	if err != nil {
		return err
	}

	handler := string(data)
	handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(arg3))
	handler = strings.ReplaceAll(handler, "$MODULENAME$", moduleName)

	err = os.WriteFile(fileName, []byte(handler), 0644)
	if err != nil {
		return err
	}

	return nil
}
