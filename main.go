package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Contact struct represents the contact object
type Contact struct {
	ID                int       `json:"id"`
	ContactName       string    `json:"contact_name"`
	PreferredChannel  []Channel `json:"preferred_channel"`
	PreferredLanguage []string  `json:"preferred_language"`
}

// Channel struct represents the preferred channel details
type Channel struct {
	ID             int    `json:"id"`
	ChannelName    string `json:"channel_name"`
	ChannelDetails string `json:"channel_details"`
}

var contacts []Contact

const dataFile = "contacts.json"

func main() {
	loadData()

	router := mux.NewRouter()

	// Create Contact API
	router.HandleFunc("/addContact", createContact).Methods("POST")

	// Get List of Contacts API
	router.HandleFunc("/getContacts", getContacts).Methods("GET")

	// Get Contact by ID API
	router.HandleFunc("/contactById/{id}", getContactByID).Methods("GET")

	router.HandleFunc("/getPreferredChannel/{id}", getPreferredChannelByID).Methods("GET")

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", router)

	// Save data when the server is stopped
	saveData()
}
func getPreferredChannelByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	for _, contact := range contacts {
		if contact.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(contact.PreferredChannel)
			return
		}
	}

	http.Error(w, "Contact not found", http.StatusNotFound)
}

func createContact(w http.ResponseWriter, r *http.Request) {
	var newContact Contact
	err := json.NewDecoder(r.Body).Decode(&newContact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newContact.ID = len(contacts) + 1
	contacts = append(contacts, newContact)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newContact)

	// Save data after creating a new contact
	saveData()
}

func getContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}

func getContactByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	for _, contact := range contacts {
		if contact.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(contact)
			return
		}
	}

	http.Error(w, "Contact not found", http.StatusNotFound)
}

func loadData() {
	file, err := os.Open(dataFile)
	if err != nil {
		// If the file doesn't exist, initialize an empty contacts list
		contacts = []Contact{}
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading data file:", err)
		return
	}

	err = json.Unmarshal(data, &contacts)
	if err != nil {
		fmt.Println("Error unmarshalling data:", err)
		return
	}
}

func saveData() {
	data, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	err = ioutil.WriteFile(dataFile, data, 0644)
	if err != nil {
		fmt.Println("Error writing data file:", err)
		return
	}
}
