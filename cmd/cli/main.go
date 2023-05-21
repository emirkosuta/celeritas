package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/emirkosuta/celeritas"
	"github.com/fatih/color"
)

const version = "1.0.0"

var moduleName string
var cel celeritas.Celeritas

func main() {
	var message string
	arg1, arg2, arg3, err := validateInput()
	if err != nil {
		exitGracefully(err)
	}

	setup(arg1)

	switch arg1 {
	case "help":
		showHelp()
	case "new":
		if arg2 == "" {
			exitGracefully(errors.New("make requires a subcommand: (migration|model|handler)"))
		}
		err = doNew(arg2)
		if err != nil {
			exitGracefully(err)
		}

	case "version":
		color.Yellow("Application version: " + version)

	case "migrate":
		if arg2 == "" {
			arg2 = "up"
		}
		err = doMigrate(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}
		message = "Migrations complete!"

	case "make":
		if arg2 == "" {
			exitGracefully(errors.New("make requires a subcommand: (migration|model|handler)"))
		}
		err = doMake(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}

	default:
		showHelp()
	}

	exitGracefully(nil, message)
}

func validateInput() (string, string, string, error) {
	var arg1, arg2, arg3 string

	if len(os.Args) > 1 {
		arg1 = os.Args[1]

		if len(os.Args) >= 3 {
			arg2 = os.Args[2]
		}

		if len(os.Args) >= 4 {
			arg3 = os.Args[3]
		}
	} else {
		showHelp()
		return "", "", "", errors.New("command required")
	}

	return arg1, arg2, arg3, nil
}

func exitGracefully(err error, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}

	if err != nil {
		showHelp()
		color.Red("Error: %v\n", err)
	}

	if len(message) > 0 {
		color.Yellow(message)
	} else {
		color.Green("Finished!")
	}

	os.Exit(0)
}

func addImportStatement(filename, importStatement string) error {
	// Read the content of the file
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Check if the import block exists in the file
	importBlockPattern := regexp.MustCompile(`import \([\s\S]*?\)`)
	importBlockMatches := importBlockPattern.FindAllIndex(content, -1)

	// If the import block exists, add the import statement inside it
	if len(importBlockMatches) > 0 {
		importBlockStart := importBlockMatches[0][0]
		content = append(content[:importBlockStart], append([]byte(importStatement+"\n"), content[importBlockStart:]...)...)
	} else {
		// If the import block doesn't exist, add a new import block
		importIndex := strings.Index(string(content), "import")
		if importIndex == -1 {
			// If there are no imports, add the import statement directly
			content = append([]byte(importStatement+"\n"), content...)
		} else {
			// If there are existing imports, insert a new import block
			content = append(content[:importIndex+len("import\n")], append([]byte(importStatement+"\n"), content[importIndex+len("import\n"):]...)...)
		}
	}

	err = os.WriteFile(filename, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil

}
