package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AniketGodambe/mongoapi/router"
)

func main() {

	fmt.Println("Mongo DB API")
	r := router.Router()

	fmt.Println("Server is getting started...")

	log.Fatal(http.ListenAndServe(":4000", r))
	log.Println("Server is running on port 4000")

}
