package router

import (
	"github.com/AniketGodambe/mongoapi/controller"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// Get List of Contacts API
	router.HandleFunc("/api/getContacts", controller.GetAllContactHandler).Methods("GET")

	router.HandleFunc("/api/addContact", controller.CreateContactHandler).Methods("POST")

	router.HandleFunc("/api/addContact", controller.UpdateContactHandler).Methods("PUT")

	router.HandleFunc("/api/delete", controller.DeleteOneContactHandler).Methods("DELETE")

	router.HandleFunc("/api/deleteAll", controller.DeleteAllContactHandler).Methods("DELETE")

	// Questions API
	router.HandleFunc("/api/questions/add", controller.AddQuestion).Methods("POST")
	router.HandleFunc("/api/questions/update", controller.UpdateQuestion).Methods("PUT")
	router.HandleFunc("/api/questions/delete", controller.DeleteQuestion).Methods("DELETE")
	router.HandleFunc("/api/questions/hide", controller.HideQuestion).Methods("PUT")

	return router

}
