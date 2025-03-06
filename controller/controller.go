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
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://aniketgodambe:aniketgodambe@cluster0.ehh0w.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
const dbName = "contactdb"
const collName = "contacts"

// most imp
var collection *mongo.Collection

// connect with mongoDB
func init() {

	//client option
	clientOptions := options.Client().ApplyURI(connectionString)
	//connect to mongo
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	collection = client.Database(dbName).Collection(collName)
}

// MONGODB HELPER \\

// Insert one record
func insertOneContact(contact model.Contact) (bool, string, int) {
	// Check if mobile number already exists
	var existingContact model.Contact
	err := collection.FindOne(context.TODO(), bson.M{"mobile": contact.Mobile}).Decode(&existingContact)
	if err == nil {
		return false, "Mobile number already exists!", 0
	} else if err != mongo.ErrNoDocuments {
		log.Println("Error checking existing contact:", err)
		return false, "Database error!", 0
	}

	// Assign a unique integer ID (auto-incremented)
	count, err := collection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		log.Println("Error counting documents:", err)
		return false, "Failed to generate user ID!", 0
	}

	contact.ID = int(count) + 1 // Simple auto-increment logic

	// Insert new contact
	result, err := collection.InsertOne(context.TODO(), contact)
	if err != nil {
		log.Println("Error inserting contact:", err)
		return false, "Failed to insert contact!", 0
	}

	fmt.Println("Inserted ID:", result.InsertedID)

	return true, "Contact inserted successfully!", contact.ID
}

// update record

func deleteOneContact(contactId string) {
	id, _ := primitive.ObjectIDFromHex(contactId)
	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Contact deleted successfully!", result.DeletedCount)
}

func deleteAllContact() {
	result, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Contact deleted successfully!", result.DeletedCount)
}

func getContactById(contactId string) model.Contact {
	var contact model.Contact
	id, _ := primitive.ObjectIDFromHex(contactId)
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.Background(), filter).Decode(&contact)
	if err != nil {
		log.Fatal(err)
	}
	return contact
}

func getAllContacts() []model.Contact {
	var contacts []model.Contact
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var contact model.Contact
		cursor.Decode(&contact)
		contacts = append(contacts, contact)
	}
	return contacts
}

func GetAllContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	contacts := getAllContacts()
	json.NewEncoder(w).Encode(contacts)
}
func CreateContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newContact model.Contact
	err := json.NewDecoder(r.Body).Decode(&newContact)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate mobile number (10 digits only)
	mobileRegex := regexp.MustCompile(`^\d{10}$`)
	if !mobileRegex.MatchString(newContact.Mobile) {
		response := map[string]interface{}{
			"success": false,
			"message": "Mobile number must be exactly 10 digits!",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	success, message, userID := insertOneContact(newContact)

	// Response JSON
	response := map[string]interface{}{
		"success": success,
		"message": message,
	}
	if success {
		response["user_id"] = userID
	}

	json.NewEncoder(w).Encode(response)
}

func DeleteOneContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := r.URL.Query()
	contactId := params.Get("id")
	deleteOneContact(contactId)
	w.Write([]byte("Contact deleted successfully!"))
}

func DeleteAllContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	deleteAllContact()
	w.Write([]byte("All contacts deleted successfully!"))
}

func updateContact(contactID int, name string, age int) (int64, error) {
	filter := bson.M{"_id": contactID}
	update := bson.M{"$set": bson.M{"contact_name": name, "age": age}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Error updating contact:", err)
		return 0, err
	}

	return result.ModifiedCount, nil
}

func UpdateContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	// Decode JSON body into Contact struct
	var contact model.Contact
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate that ID is present
	if contact.ID == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Call update function
	updatedCount, err := updateContact(contact.ID, contact.ContactName, contact.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	response := fmt.Sprintf("Contact with ID %d updated successfully!", contact.ID)
	if updatedCount == 0 {
		response = "No contact was updated, possibly wrong ID."
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
