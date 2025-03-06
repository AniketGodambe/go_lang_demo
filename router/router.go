package router

import (
	"github.com/AniketGodambe/mongoapi/controller"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// Get List of Contacts API
	router.HandleFunc("/api/getContacts", controller.DeleteOneContactHandler).Methods("GET")

	router.HandleFunc("/api/addContact", controller.CreateContactHandler).Methods("POST")

	router.HandleFunc("/api/addContact", controller.UpdateContactHandler).Methods("PUT")

	router.HandleFunc("/api/delete", controller.DeleteOneContactHandler).Methods("DELETE")

	router.HandleFunc("/api/deleteAll", controller.DeleteAllContactHandler).Methods("DELETE")

	return router

}
