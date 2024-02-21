package main

import (
	"log"

	"github.com/namatoj/sociallink/pkg/web"
)

func main() {
	app := web.App()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
