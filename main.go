package main

import (
	"final-project/config/postgres"
	"final-project/routes"

	"log"
	"os"

	"github.com/subosito/gotenv"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if err := gotenv.Load(); err != nil {
		log.Println(err)
	}
}

func main() {
	postgres.Connect()
	router := routes.StartApp()
	router.Run(":" + os.Getenv("PORT"))
}
