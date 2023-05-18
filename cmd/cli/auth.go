package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

func doAuth() error {
	// migrations
	dbType := cel.DB.DatabaseType
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	upFile := cel.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := cel.RootPath + "/migrations/" + fileName + ".down.sql"

	err := copyFilefromTemplate("templates/migrations/auth_tables."+dbType+".sql", upFile)
	if err != nil {
		return err
	}

	err = copyDataToFile([]byte("drop table if exists users cascade;"), downFile)
	if err != nil {
		return err
	}

	// copy over auth handler
	userHandlerFile, err := templateFS.ReadFile("templates/handlers/user.go.txt")
	if err != nil {
		return err
	}
	userHandler := string(userHandlerFile)
	userHandler = strings.ReplaceAll(userHandler, "$MODULENAME$", moduleName)
	err = copyDataToFile([]byte(userHandler), cel.RootPath+"/handlers/user.go")
	if err != nil {
		return err
	}

	//copy over user service
	userServiceFile, err := templateFS.ReadFile("templates/services/user.go.txt")
	if err != nil {
		return err
	}
	userService := string(userServiceFile)
	userService = strings.ReplaceAll(userService, "$MODULENAME$", moduleName)
	err = copyDataToFile([]byte(userService), cel.RootPath+"/services/user.go")
	if err != nil {
		return err
	}

	// copy over user model and register it in models
	err = copyFilefromTemplate("templates/data/user.go.txt", cel.RootPath+"/data/user.go")
	if err != nil {
		return err
	}
	err = registerModel("User")
	if err != nil {
		return err
	}

	//copy user dto
	userDtoFile, err := templateFS.ReadFile("templates/dto/user.go.txt")
	if err != nil {
		return err
	}
	userDto := string(userDtoFile)
	userDto = strings.ReplaceAll(userDto, "$MODULENAME$", moduleName)
	err = copyDataToFile([]byte(userDto), cel.RootPath+"/dto/user.go")
	if err != nil {
		return err
	}

	// copy over auth middleware
	authMiddlewareFile, err := templateFS.ReadFile("templates/middleware/auth.go.txt")
	if err != nil {
		return err
	}
	authMiddleware := string(authMiddlewareFile)
	authMiddleware = strings.ReplaceAll(authMiddleware, "$MODULENAME$", moduleName)
	err = copyDataToFile([]byte(authMiddleware), cel.RootPath+"/middleware/auth.go")
	if err != nil {
		return err
	}

	// copy over auth handler
	authHandlerFile, err := templateFS.ReadFile("templates/handlers/auth-handlers.go.txt")
	if err != nil {
		return err
	}
	authHandler := string(authHandlerFile)
	authHandler = strings.ReplaceAll(authHandler, "$MODULENAME$", moduleName)
	err = copyDataToFile([]byte(authHandler), cel.RootPath+"/handlers/auth-handlers.go")
	if err != nil {
		return err
	}

	// copy over auth handler
	authRoutes, err := templateFS.ReadFile("templates/routes/auth.go.txt")
	if err != nil {
		return err
	}
	err = addApiRoute(string(authRoutes))
	if err != nil {
		return err
	}

	//copy over auth dto
	err = copyFilefromTemplate("templates/dto/auth.go.txt", cel.RootPath+"/dto/auth.go")
	if err != nil {
		return err
	}

	//copy over auth service
	authServiceFile, err := templateFS.ReadFile("templates/services/auth.go.txt")
	if err != nil {
		return err
	}
	authService := string(authServiceFile)
	authService = strings.ReplaceAll(authService, "$MODULENAME$", moduleName)
	err = copyDataToFile([]byte(authService), cel.RootPath+"/services/auth.go")
	if err != nil {
		return err
	}

	err = copyFilefromTemplate("templates/mailer/password-reset.html.tmpl", cel.RootPath+"/mail/password-reset.html.tmpl")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/mailer/password-reset.plain.tmpl", cel.RootPath+"/mail/password-reset.plain.tmpl")
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("  - users migration created. Don't forget to execute the migration.")
	color.Yellow("  - user model, dto, service and handler created")
	color.Yellow("  - password reset mail created")
	color.Yellow("  - auth middleware created")
	color.Yellow("  - auth handler, service and dto created")

	return nil
}
