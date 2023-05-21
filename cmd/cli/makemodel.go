package main

import (
	"errors"
	"os"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

func doModel(arg3 string) error {
	if arg3 == "" {
		return errors.New("you must give the model a name")
	}

	data, err := templateFS.ReadFile("templates/data/model.go.txt")
	if err != nil {
		return err
	}

	model := string(data)

	plur := pluralize.NewClient()

	var modelName = arg3
	var tableName = arg3

	if plur.IsPlural(arg3) {
		modelName = plur.Singular(arg3)
		tableName = strings.ToLower(tableName)
	} else {
		tableName = strings.ToLower(plur.Plural(arg3))
	}
	modelName = strcase.ToCamel(modelName)

	fileName := cel.RootPath + "/data/" + strings.ToLower(modelName) + ".go"
	if fileExists(fileName) {
		return err
	}

	model = strings.ReplaceAll(model, "$MODELNAME$", modelName)
	model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

	err = copyDataToFile([]byte(model), fileName)
	if err != nil {
		return err
	}
	err = makeMigration(tableName)
	if err != nil {
		return err
	}

	// register model
	err = registerModel(modelName)
	if err != nil {
		return err
	}

	return nil
}

func findClosingBraceIndex(content string, fromIndex int) (int, error) {
	openBraces := 0
	for i := fromIndex; i < len(content); i++ {
		switch content[i] {
		case '{':
			openBraces++
		case '}':
			openBraces--
			if openBraces == 0 {
				return i, nil
			}
		}
	}
	return -1, errors.New("matching closing brace not found")
}

func registerModel(modelName string) error {
	modelsdata, err := os.ReadFile(cel.RootPath + "/data/models.go")
	if err != nil {
		return err
	}
	modelsContent := string(modelsdata)

	// Find the insertion point
	insertIndex, err := findSubstringIndex(modelsContent, "type Models struct", 0)
	if err != nil {
		return errors.New("'type Models struct' not found")
	}

	// Find the next closing curly brace after the insertion point
	registerModelPoint1, err := findSubstringIndex(modelsContent, "}", insertIndex)
	if err != nil {
		return errors.New("'Register model point not found")
	}

	// Insert your text on a new line before the closing brace
	modelsContent = modelsContent[:registerModelPoint1] + "\t" + modelName + " " + modelName + "\n\t" + modelsContent[registerModelPoint1:]

	// Find the insertion point
	insertIndex2, err := findSubstringIndex(modelsContent, "return Models", 0)
	if err != nil {
		return errors.New("'return Models' not found")
	}

	// Find the next closing curly brace after the insertion point
	registerModelPoint2, err := findClosingBraceIndex(modelsContent, insertIndex2)
	if err != nil {
		return errors.New("'Register model point not found")
	}

	// Insert your text on a new line before the closing brace
	modelsContent = modelsContent[:registerModelPoint2] + "\t" + modelName + ": " + modelName + "{},\n\t" + modelsContent[registerModelPoint2:]

	err = copyDataToFile([]byte(modelsContent), cel.RootPath+"/data/models.go")
	if err != nil {
		return err
	}

	return nil
}
