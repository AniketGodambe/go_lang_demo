package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AniketGodambe/mongoapi/controller"
	"github.com/AniketGodambe/mongoapi/router"
)

func main() {

	controller.InitDB()

	fmt.Println("Mongo DB API")
	r := router.Router()

	fmt.Println("Server is getting started...")

	log.Fatal(http.ListenAndServe(":8080", r))
	log.Println("Server is running on port 8080")

}
