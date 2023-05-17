package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gertd/go-pluralize"
)

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
