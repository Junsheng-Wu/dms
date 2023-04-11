package main

import (
	"log"

	"dms/cmd/dms/app"
)

func main() {
	cmd := app.NewAPIServerCommand()

	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
