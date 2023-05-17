package main

import (
	"errors"
	"strings"

	"github.com/fatih/color"
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
	case "auth-api":
		err := doAuthApi()
		if err != nil {
			exitGracefully(err)
		}
	case "handler":
		err := makeHandler(arg3)
		if err != nil {
			exitGracefully(err)
		}
		color.Yellow("Created the handler.")
	case "service":
		err := makeService(arg3)
		if err != nil {
			exitGracefully(err)
		}
		color.Yellow("Created the service and dto.")
	case "model":
		err := doModel(arg3)
		if err != nil {
			exitGracefully(err)
		}
		color.Yellow("Created the model and migration.")
	case "resource":
		err := doModel(arg3)
		if err != nil {
			exitGracefully(err)
		}
		err = makeService(arg3)
		if err != nil {
			exitGracefully(err)
		}
		err = makeHandler(arg3)
		if err != nil {
			exitGracefully(err)
		}
		color.Yellow("Created the model, migration, service and handler for " + arg3 + " entity.")
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

func findSubstringIndex(content, substring string, fromIndex int) (int, error) {
	index := strings.Index(content[fromIndex:], substring)
	if index == -1 {
		return -1, errors.New("'" + substring + "' not found")
	}
	return index + fromIndex, nil
}
