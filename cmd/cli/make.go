package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

func doMake(arg2, arg3 string) error {

	switch arg2 {
	case "key":
		rnd := cel.RandomString(32)
		color.Yellow("32 character encyption key: %s", rnd)
	case "migration":
		err := makeMigration(arg3)
		if err != nil {
			exitGracefully(err)
		}
	case "auth":
		err := doAuth()
		if err != nil {
			exitGracefully(err)
		}
	case "handler":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the handler a name"))
		}

		fileName := cel.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exists!"))
		}

		data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		handler := string(data)
		handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(arg3))

		err = os.WriteFile(fileName, []byte(handler), 0644)
		if err != nil {
			exitGracefully(err)
		}
	case "model":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the model a name"))
		}

		data, err := templateFS.ReadFile("templates/data/model.go.txt")
		if err != nil {
			exitGracefully(err)
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

		fileName := cel.RootPath + "/data/" + strings.ToLower(modelName) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exists!"))
		}

		model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
		model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

		err = copyDataToFile([]byte(model), fileName)
		if err != nil {
			exitGracefully(err)
		}
		err = makeMigration(tableName)
		if err != nil {
			exitGracefully(err)
		}
		color.Yellow("Created the model and migration.")
	case "session":
		err := doSessionTable()
		if err != nil {
			exitGracefully(err)
		}
	case "mail":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the mail template a name"))
		}
		htmlMail := cel.RootPath + "/mail/" + strings.ToLower(arg3) + ".html.tmpl"
		plainMail := cel.RootPath + "/mail/" + strings.ToLower(arg3) + ".plain.tmpl"

		err := copyFilefromTemplate("templates/mailer/mail.html.tmpl", htmlMail)
		if err != nil {
			exitGracefully(err)
		}

		err = copyFilefromTemplate("templates/mailer/mail.plain.tmpl", plainMail)
		if err != nil {
			exitGracefully(err)
		}
	default:
		return errors.New("make " + arg2 + " is not supported.")
	}

	return nil
}

func makeMigration(name string) error {
	dbType := cel.DB.DatabaseType

	if name == "" {
		return errors.New("you must give the migration a name")
	}

	plur := pluralize.NewClient()

	if plur.IsPlural(name) {
		name = strings.ToLower(name)
	} else {
		name = strings.ToLower(plur.Plural(name))
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), name)
	upFileName := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
	downFileName := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

	upFileData, err := templateFS.ReadFile("templates/migrations/migration." + dbType + ".up.sql")
	if err != nil {
		return err
	}

	migrationUp := string(upFileData)
	migrationUp = strings.ReplaceAll(migrationUp, "$MIGRATIONNAME$", name)

	err = os.WriteFile(upFileName, []byte(migrationUp), 0644)
	if err != nil {
		return err
	}

	data, err := templateFS.ReadFile("templates/migrations/migration." + dbType + ".down.sql")
	if err != nil {
		return err
	}

	migrationDown := string(data)
	migrationDown = strings.ReplaceAll(migrationDown, "$MIGRATIONNAME$", name)

	err = os.WriteFile(downFileName, []byte(migrationDown), 0644)
	if err != nil {
		return err
	}
	return nil
}
