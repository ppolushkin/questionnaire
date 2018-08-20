package main

import (
	"fmt"
	"os"
)

func main() {
	a := App{}
	appPort := env("appPort", ":8080")
	dbHost := env("dbHost", "localhost")
	dbName := env("dbName", "quest")
	dbUser := env("dbUser", "quest")
	dbPswd := env("dbPswd", "quest")

	const startingMsg = "Starting questionnaire application on \n DB host: %s\n DB name: %s\n"
	fmt.Printf(startingMsg, dbHost, dbName)

	if err := a.Initialize(dbUser, dbPswd, dbHost, dbName); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Initialized successfully!\nStarted on localhost%d", appPort)
	a.Run(appPort)
}

func env(key, defaultVal string) string {
	appPort, ok := os.LookupEnv(key)
	if !ok {
		appPort = defaultVal
	}
	return appPort
}
