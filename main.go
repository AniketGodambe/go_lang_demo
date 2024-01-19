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

type User struct {
	UserId   int    `json:"userId"`
	UserName string `json:"user_name"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	City     string `json:"city"`
	Age      int    `json:"age"`
}

// Channel struct represents the preferred channel details
type Channel struct {
	ID             int    `json:"id"`
	ChannelName    string `json:"channel_name"`
	ChannelDetails string `json:"channel_details"`
}

var contacts []Contact
var users []User

const contacstsFile = "contacts.json"
const usersFile = "users.json"

func main() {
	loadContactsData()
	loadUsersData()
	router := mux.NewRouter()

	// Create Contact API
	router.HandleFunc("/addContact", createContact).Methods("POST")

	// Get List of Contacts API
	router.HandleFunc("/getContacts", getContacts).Methods("GET")

	// Get Contact by ID API
	router.HandleFunc("/contactById/{id}", getContactByID).Methods("GET")
	router.HandleFunc("/getPreferredChannel/{id}", getPreferredChannelByID).Methods("GET")

	//Create User API
	router.HandleFunc("/createUser", createUser).Methods("POST")
	// Get List of Users API
	router.HandleFunc("/getUserList", getUsers).Methods("GET")

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", router)

	// Save data when the server is stopped
	saveContactsData()
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
	saveContactsData()
}
func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the email already exists
	if emailExists(newUser.Email) {
		response := map[string]interface{}{
			"success": false,
			"message": "Email already exists",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	if phoneExists(newUser.Mobile) {
		response := map[string]interface{}{
			"success": false,
			"message": "Mobile number already exists",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	newUser.UserId = len(users) + 1
	users = append(users, newUser)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"data":    newUser,
		"message": "User created successfully",
	}
	json.NewEncoder(w).Encode(response)

	saveUsersData()
}

// Function to check if the email already exists
func emailExists(email string) bool {
	for _, user := range users {
		if user.Email == email {
			return true
		}
	}
	return false
}

func phoneExists(mobile string) bool {
	for _, user := range users {
		if user.Mobile == mobile {
			return true
		}
	}
	return false
}

func getContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
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

func loadContactsData() {
	file, err := os.Open(contacstsFile)
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

func saveContactsData() {
	data, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	err = ioutil.WriteFile(contacstsFile, data, 0644)
	if err != nil {
		fmt.Println("Error writing data file:", err)
		return
	}
}

func loadUsersData() {
	file, err := os.Open(usersFile)
	if err != nil {
		// If the file doesn't exist, initialize an empty contacts list
		users = []User{}
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading data file:", err)
		return
	}

	err = json.Unmarshal(data, &users)
	if err != nil {
		fmt.Println("Error unmarshalling data:", err)
		return
	}
}

func saveUsersData() {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	err = ioutil.WriteFile(usersFile, data, 0644)
	if err != nil {
		fmt.Println("Error writing data file:", err)
		return
	}
}
