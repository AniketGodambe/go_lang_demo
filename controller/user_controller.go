package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/AniketGodambe/mongoapi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// getAllContacts fetches all contacts from the database
func getAllContacts() ([]model.Contact, error) {
	var contacts []model.Contact
	cursor, err := ContactsCollection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var contact model.Contact
		if err := cursor.Decode(&contact); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// GetAllContactHandler handles the API request to fetch all contacts
func GetAllContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	contacts, err := getAllContacts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Failed to retrieve contacts",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{
		Message:    "Success",
		StatusCode: http.StatusOK,
		Data:       contacts,
	})
}

func createOneContact(contact model.Contact) (bool, string, int) {
	var existingContact model.Contact
	err := ContactsCollection.FindOne(context.TODO(), bson.M{"mobile": contact.Mobile}).Decode(&existingContact)
	if err == nil {
		return false, "Mobile number already exists!", 0
	} else if err != mongo.ErrNoDocuments {
		log.Println("Error checking existing contact:", err)
		return false, "Database error!", 0
	}

	count, err := ContactsCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		log.Println("Error counting documents:", err)
		return false, "Failed to generate user ID!", 0
	}

	contact.ID = int(count) + 1

	result, err := ContactsCollection.InsertOne(context.TODO(), contact)
	if err != nil {
		log.Println("Error inserting contact:", err)
		return false, "Failed to insert contact!", 0
	}

	fmt.Println("Inserted ID:", result.InsertedID)

	return true, "Contact inserted successfully!", contact.ID
}

// CreateContactHandler handles API request to add a new contact
func CreateContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Invalid request method",
			StatusCode: http.StatusMethodNotAllowed,
		})
		return
	}

	var newContact model.Contact
	err := json.NewDecoder(r.Body).Decode(&newContact)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Invalid request body",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// Validate mobile number (must be exactly 10 digits)
	mobileRegex := regexp.MustCompile(`^\d{10}$`)
	if !mobileRegex.MatchString(newContact.Mobile) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Mobile number must be exactly 10 digits!",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	success, message, userID := createOneContact(newContact)

	statusCode := http.StatusCreated
	if !success {
		statusCode = http.StatusConflict
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.Response{
		Message:    message,
		StatusCode: statusCode,
		Data:       map[string]int{"user_id": userID},
	})
}

// deleteOneContact removes a contact by ID from the database
func deleteOneContact(contactId string) (bool, string, int64) {
	id, err := primitive.ObjectIDFromHex(contactId)
	if err != nil {
		log.Println("Invalid contact ID format:", err)
		return false, "Invalid contact ID format!", 0
	}

	filter := bson.M{"_id": id}
	result, err := ContactsCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println("Error deleting contact:", err)
		return false, "Database error!", 0
	}

	if result.DeletedCount == 0 {
		return false, "Contact not found!", 0
	}

	fmt.Println("Contact deleted successfully!", result.DeletedCount)
	return true, "Contact deleted successfully!", result.DeletedCount
}

// DeleteOneContactHandler handles API requests to delete a contact
func DeleteOneContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Invalid request method",
			StatusCode: http.StatusMethodNotAllowed,
		})
		return
	}

	params := r.URL.Query()
	contactId := params.Get("id")
	if contactId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Missing contact ID",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	success, message, deletedCount := deleteOneContact(contactId)

	statusCode := http.StatusOK
	if !success {
		if deletedCount == 0 {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.Response{
		Message:    message,
		StatusCode: statusCode,
		Data:       map[string]int64{"deleted_count": deletedCount},
	})
}

func deleteAllContact() (int64, error) {
	result, err := ContactsCollection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func DeleteAllContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	deletedCount, err := deleteAllContact()
	if err != nil {
		response := model.Response{
			Message:    "Failed to delete contacts",
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := model.Response{
		Message:    "All contacts deleted successfully!",
		StatusCode: http.StatusOK,
		Data:       map[string]int64{"deletedCount": deletedCount},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	var contact model.Contact
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if contact.ID == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	updatedCount, err := updateContact(contact.ID, contact.ContactName, contact.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("Contact with ID %d updated successfully!", contact.ID)
	if updatedCount == 0 {
		response = "No contact was updated, possibly wrong ID."
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func updateContact(contactID int, name string, age int) (int64, error) {
	filter := bson.M{"_id": contactID}
	update := bson.M{"$set": bson.M{"contact_name": name, "age": age}}

	result, err := ContactsCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Error updating contact:", err)
		return 0, err
	}

	return result.ModifiedCount, nil
}
