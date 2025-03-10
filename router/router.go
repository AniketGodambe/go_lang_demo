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
	router.HandleFunc("/api/questions/add", controller.AddQuestionHandler).Methods("POST")
	router.HandleFunc("/api/questions/update", controller.UpdateQuestionHandler).Methods("PUT")
	router.HandleFunc("/api/questions/delete", controller.DeleteQuestionHandler).Methods("DELETE")
	router.HandleFunc("/api/questions/questionVisibility", controller.ToggleQuestionVisibilityHandler).Methods("PUT")
	router.HandleFunc("/api/questions/questionsList", controller.GetAllQuestionsHandler).Methods("GET")

	router.HandleFunc("/api/getQuestionById", controller.GetQuestionByIdHandler).Methods("GET")

	return router

}
