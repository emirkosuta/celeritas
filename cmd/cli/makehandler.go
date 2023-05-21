package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

func makeHandler(arg3 string) error {
	if arg3 == "" {
		return errors.New("you must give the handler a name")
	}

	plur := pluralize.NewClient()

	var handlerName = arg3
	var routeBaseName = arg3

	if plur.IsPlural(arg3) {
		handlerName = strcase.ToCamel(plur.Singular(arg3))
		routeBaseName = strings.ToLower(routeBaseName)
	} else {
		handlerName = strcase.ToCamel(arg3)
		routeBaseName = strings.ToLower(plur.Plural(arg3))
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
	handler = strings.ReplaceAll(handler, "$HANDLERNAME$", handlerName)
	handler = strings.ReplaceAll(handler, "$LOWERCASEHANDLER$", strings.ToLower(handlerName))
	handler = strings.ReplaceAll(handler, "$MODULENAME$", moduleName)

	err = os.WriteFile(fileName, []byte(handler), 0644)
	if err != nil {
		return err
	}

	handlerInterfaceData, err := buildHandlerInterfaceData(handlerName)
	if err != nil {
		return err
	}

	err = insertHandlerInterface(handlerName, handlerInterfaceData)
	if err != nil {
		return err
	}

	err = wireServiceAndHandler(handlerName, handlerName)
	if err != nil {
		return err
	}

	routes, err := generateResourceRoutes(handlerName, routeBaseName)
	if err != nil {
		return err
	}
	err = addApiRoute(routes)
	if err != nil {
		return err
	}

	return nil
}

func buildHandlerInterfaceData(handlerName string) (string, error) {
	handlerInterface, err := templateFS.ReadFile("templates/handlers/handler-interface.go.txt")
	if err != nil {
		return "", err
	}
	handlerInterfaceData := strings.ReplaceAll(string(handlerInterface), "$HANDLERNAME$", handlerName)

	return handlerInterfaceData, nil
}

func insertHandlerInterface(handlerName, handlerInterfaceData string) error {
	handlerdata, err := os.ReadFile(cel.RootPath + "/handlers/handlers.go")
	if err != nil {
		return err
	}
	handlerContent := string(handlerdata)

	handlerContent += handlerInterfaceData

	// Find the insertion point
	insertIndex, err := findSubstringIndex(handlerContent, "type Handlers struct", 0)
	if err != nil {
		return errors.New("'type Handlers struct' not found")
	}

	// Find the next closing curly brace after the insertion point
	registerHandlerPoint, err := findClosingBraceIndex(handlerContent, insertIndex)
	if err != nil {
		return errors.New("'Register handler point not found")
	}

	handlerContent = handlerContent[:registerHandlerPoint] + "\t" + handlerName + "Handler " + handlerName + "Handler\n\t" + handlerContent[registerHandlerPoint:]

	err = copyDataToFile([]byte(handlerContent), cel.RootPath+"/handlers/handlers.go")
	if err != nil {
		return err
	}

	err = addImportStatement(cel.RootPath+"/handlers/handlers.go", `"net/http"`)
	if err != nil {
		return err
	}

	return nil
}

func wireServiceAndHandler(handlerName string, modelName string) error {
	initAppData, err := os.ReadFile(cel.RootPath + "/init-app.go")
	if err != nil {
		return err
	}
	initAppContent := string(initAppData)

	// Find the insertion point
	insertIndex, err := findSubstringIndex(initAppContent, "myHandlers := &handlers.Handlers", 0)
	if err != nil {
		return errors.New("'return app' not found")
	}

	wireServiceContent := fmt.Sprintf(`
	%sService := services.New%sServiceImpl(cel, models.%s)
	%sHandler := handlers.New%sHandler(cel, %sService)`,
		strings.ToLower(handlerName),
		handlerName,
		modelName,
		strings.ToLower(handlerName),
		handlerName,
		strings.ToLower(handlerName),
	)

	initAppContent = initAppContent[:insertIndex] + "\t\n" + wireServiceContent + "\n\n\t" + initAppContent[insertIndex:]

	// Find the next closing curly brace after the insertion point
	registerHandlerPoint, err := findClosingBraceIndex(initAppContent, insertIndex)
	if err != nil {
		return errors.New("'Register model point not found")
	}

	initAppContent = initAppContent[:registerHandlerPoint] + "\t" + handlerName + "Handler: " + strings.ToLower(handlerName) + "Handler,\n\t" + initAppContent[registerHandlerPoint:]

	err = copyDataToFile([]byte(initAppContent), cel.RootPath+"/init-app.go")
	if err != nil {
		return err
	}

	err = addImportStatement(cel.RootPath+"/init-app.go", `"`+moduleName+`/services"`)
	if err != nil {
		return err
	}

	return nil
}

func generateResourceRoutes(entity, routeBaseName string) (string, error) {
	data, err := templateFS.ReadFile("templates/routes/resource.go.txt")
	if err != nil {
		return "", err
	}

	routes := string(data)
	routes = strings.ReplaceAll(routes, "$MODELNAME$", entity)
	routes = strings.ReplaceAll(routes, "$ROUTEBASENAME$", routeBaseName)

	return routes, nil
}

func addApiRoute(routes string) error {
	routedata, err := os.ReadFile(cel.RootPath + "/routes.go")
	if err != nil {
		return err
	}
	routeContent := string(routedata)

	// Find the insertion point
	insertIndex, err := findSubstringIndex(routeContent, `Route("/api"`, 0)
	if err != nil {
		return errors.New(`Route("/api" not found`)
	}

	// Find the next closing curly brace after the insertion point
	registerRoutePoint, err := findClosingBraceIndex(routeContent, insertIndex)
	if err != nil {
		return errors.New("register model point not found")
	}

	// Insert your text on a new line before the closing brace
	routeContent = routeContent[:registerRoutePoint] + "\n\t\t" + routes + "\n\t" + routeContent[registerRoutePoint:]

	err = copyDataToFile([]byte(routeContent), cel.RootPath+"/routes.go")
	if err != nil {
		return err
	}

	return nil
}
