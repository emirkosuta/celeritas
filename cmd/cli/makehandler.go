package main

import (
	"errors"
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
	handler = strings.ReplaceAll(handler, "$MODULENAME$", moduleName)

	err = os.WriteFile(fileName, []byte(handler), 0644)
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
