package main

import "errors"

func doMigrate(arg2, arg3 string) error {
	dsn := getDSN()

	// run the migration command
	switch arg2 {
	case "up":
		err := cel.MigrateUp(dsn)
		if err != nil {
			return err
		}

	case "down":
		if arg3 == "all" {
			err := cel.MigrateDownAll(dsn)
			if err != nil {
				return err
			}
		} else {
			err := cel.Steps(-1, dsn)
			if err != nil {
				return err
			}
		}
	case "reset":
		err := cel.MigrateDownAll(dsn)
		if err != nil {
			return err
		}
		err = cel.MigrateUp(dsn)
		if err != nil {
			return err
		}
	case "force":
		err := cel.MigrateForce(dsn)
		if err != nil {
			return err
		}
	default:
		return errors.New(`Arg "` + arg2 + `" is not supported.`)
	}

	return nil
}
