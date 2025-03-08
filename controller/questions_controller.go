package controller

import (
	"context"
	"encoding/json"
	"strconv"

	"log"
	"net/http"
	"time"

	"github.com/AniketGodambe/mongoapi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getAllQuestions() ([]model.Question, error) {
	var questions []model.Question
	cursor, err := QuestionsCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &questions); err != nil {
		return nil, err
	}
	return questions, nil
}

// GetAllQuestionsHandler handles API request to fetch all questions
func GetAllQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, http.MethodGet)
	questions, err := getAllQuestions()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve questions")
		return
	}
	respondWithJSON(w, http.StatusOK, questions)
}

// Create a new question
func createOneQuestion(question model.Question) (bool, string, int) {
	// Check if the question already exists
	existingQuestion := QuestionsCollection.FindOne(context.TODO(), bson.M{"question": question.Question})
	if existingQuestion.Err() == nil {
		return false, "Question already exists!", 0
	}

	// Generate a new ID
	count, err := QuestionsCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return false, "Failed to generate question ID!", 0
	}

	question.ID = int(count) + 1
	question.CreatedAt = time.Now()
	question.LastModified = time.Now()

	// Insert the new question
	_, err = QuestionsCollection.InsertOne(context.TODO(), question)
	if err != nil {
		return false, "Failed to insert question!", 0
	}

	return true, "Question inserted successfully!", question.ID
}

// AddQuestionHandler handles API request to add a new question
func AddQuestionHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, http.MethodPost)

	var newQuestion model.Question
	if err := json.NewDecoder(r.Body).Decode(&newQuestion); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	success, _, questionID := createOneQuestion(newQuestion)
	statusCode := http.StatusCreated
	if !success {
		statusCode = http.StatusConflict
	}

	respondWithJSON(w, statusCode, map[string]int{"question_id": questionID})
}

func updateQuestion(updatedQuestion model.Question) (bool, string) {
	// Check if the question exists
	var existingQuestion model.Question
	err := QuestionsCollection.FindOne(context.TODO(), bson.M{"id": updatedQuestion.ID}).Decode(&existingQuestion)
	if err != nil {
		return false, "Question not found!"
	}

	// Check if the new question text already exists (excluding the current question)
	filter := bson.M{
		"question": updatedQuestion.Question,
		"id":       bson.M{"$ne": updatedQuestion.ID}, // Ensures the same question isn't being duplicated
	}

	count, err := QuestionsCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Println("Error checking for duplicate questions:", err)
		return false, "Failed to validate question uniqueness!"
	}

	if count > 0 {
		return false, "A question with this text already exists!"
	}

	// Update fields
	update := bson.M{
		"$set": bson.M{
			"question":       updatedQuestion.Question,
			"options":        updatedQuestion.Options,
			"correct_answer": updatedQuestion.CorrectAns,
			"reason":         updatedQuestion.Reason,
			"hidden":         updatedQuestion.Hidden,
			"last_modified":  time.Now(),
		},
	}

	// Perform update operation
	_, err = QuestionsCollection.UpdateOne(context.TODO(), bson.M{"id": updatedQuestion.ID}, update)
	if err != nil {
		log.Println("Error updating question:", err)
		return false, "Failed to update question!"
	}

	return true, "Question updated successfully!"
}

// Update an existing question
func UpdateQuestionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Invalid request method",
			StatusCode: http.StatusMethodNotAllowed,
		})
		return
	}

	// Parse request body
	var updatedQuestion model.Question
	err := json.NewDecoder(r.Body).Decode(&updatedQuestion)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Invalid request body",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// Validate ID
	if updatedQuestion.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Question ID is required",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// Call function to update question
	success, message := updateQuestion(updatedQuestion)

	statusCode := http.StatusOK
	if !success {
		statusCode = http.StatusConflict
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.Response{
		Message:    message,
		StatusCode: statusCode,
	})
}

// Delete a question
func DeleteQuestionHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Invalid question ID",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": id}

	result, err := QuestionsCollection.DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Question not found",
			StatusCode: http.StatusNotFound,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{
		Message:    "Question deleted successfully",
		StatusCode: http.StatusOK,
	})
}

// Toggle hide/show question
func toggleQuestionVisibility(questionID int) (bool, string, bool) {
	// Find the existing question
	var question model.Question
	err := QuestionsCollection.FindOne(context.TODO(), bson.M{"id": questionID}).Decode(&question)
	if err != nil {
		return false, "Question not found!", false
	}

	// Toggle the `hidden` status
	newHiddenStatus := !question.Hidden

	// Update the question in the database
	update := bson.M{
		"$set": bson.M{
			"hidden":        newHiddenStatus,
			"last_modified": time.Now(),
		},
	}

	_, err = QuestionsCollection.UpdateOne(context.TODO(), bson.M{"id": questionID}, update)
	if err != nil {
		log.Println("Error updating question visibility:", err)
		return false, "Failed to toggle question visibility!", question.Hidden
	}

	// Return success and the new status
	statusMessage := "Question is now hidden"
	if !newHiddenStatus {
		statusMessage = "Question is now visible"
	}

	return true, statusMessage, newHiddenStatus
}

func ToggleQuestionVisibilityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Invalid request method",
			StatusCode: http.StatusMethodNotAllowed,
		})
		return
	}

	// Parse request body
	var request struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Invalid request body",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// Validate ID
	if request.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Question ID is required",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// Call toggle function
	success, message, newStatus := toggleQuestionVisibility(request.ID)

	statusCode := http.StatusOK
	if !success {
		statusCode = http.StatusNotFound
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.Response{
		Message:    message,
		StatusCode: statusCode,
		Data:       map[string]bool{"hidden": newStatus},
	})
}

func getQuestionById(id int) (*model.Question, error) {
	var question model.Question
	filter := bson.M{"id": id}

	err := QuestionsCollection.FindOne(context.Background(), filter).Decode(&question)
	if err != nil {
		return nil, err
	}

	return &question, nil
}

func GetQuestionByIdHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, http.MethodGet)

	// Parse the query parameter "id" from the URL
	queryValues := r.URL.Query()
	idStr := queryValues.Get("id")

	// Validate ID parameter
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "Missing id parameter")
		return
	}

	// Convert idStr to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	// Fetch the question from database
	question, err := getQuestionById(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Question not found")
		return
	}

	// Send response
	respondWithJSON(w, http.StatusOK, question)
}

// Utility Functions
func setHeaders(w http.ResponseWriter, method string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", method)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.Response{
		Message:    message,
		StatusCode: statusCode,
		Data:       "Error",
	})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.Response{
		Message:    "Success",
		StatusCode: statusCode,
		Data:       data,
	})
}
