package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

func main() {
	app := pocketbase.New()

	app.RootCmd.AddCommand(&cobra.Command{
		Use: "hello",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Hello world!")
		},
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
