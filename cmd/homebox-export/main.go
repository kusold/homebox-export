package main

import (
	"log"
	"os"

	"github.com/kusold/homebox-export/cmd/cli"
)

func main() {
	app := cli.New()
	if err := app.Execute(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
