package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AniketGodambe/mongoapi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Fetch all questions from MongoDB
func getAllQuestions() ([]model.Question, error) {
	var questions []model.Question
	cursor, err := QuestionsCollection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var question model.Question
		if err := cursor.Decode(&question); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil
}

// GetAllQuestionsHandler handles API request to fetch all questions
func GetAllQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	questions, err := getAllQuestions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Failed to retrieve questions",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{
		Message:    "Success",
		StatusCode: http.StatusOK,
		Data:       questions,
	})
}

// Add a new question
func createOneQuestion(question model.Question) (bool, string, int) {
	// Count the existing questions to generate a new serial ID
	count, err := QuestionsCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		log.Println("Error counting documents:", err)
		return false, "Failed to generate question ID!", 0
	}

	question.ID = int(count) + 1 // Assign a new serial ID (based on document count)

	// Insert the question into the collection
	result, err := QuestionsCollection.InsertOne(context.TODO(), question)
	if err != nil {
		log.Println("Error inserting question:", err)
		return false, "Failed to insert question!", 0
	}

	fmt.Println("Inserted ID:", result.InsertedID)

	return true, "Question inserted successfully!", question.ID
}

// AddQuestionHandler handles API request to add a new question
func AddQuestionHandler(w http.ResponseWriter, r *http.Request) {
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

	var newQuestion model.Question
	err := json.NewDecoder(r.Body).Decode(&newQuestion)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Invalid request body",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// Generate and insert the question
	success, message, questionID := createOneQuestion(newQuestion)

	statusCode := http.StatusCreated
	if !success {
		statusCode = http.StatusConflict
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.Response{
		Message:    message,
		StatusCode: statusCode,
		Data:       map[string]int{"question_id": questionID},
	})
}

// Update an existing question
func UpdateQuestionHandler(w http.ResponseWriter, r *http.Request) {
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

	var q model.Question
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Invalid request",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": id}
	update := bson.M{"$set": q}

	result, err := QuestionsCollection.UpdateOne(ctx, filter, update)
	if err != nil || result.ModifiedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Question not found or not updated",
			StatusCode: http.StatusNotFound,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{
		Message:    "Question updated successfully",
		StatusCode: http.StatusOK,
		Data:       q,
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
func HideQuestionHandler(w http.ResponseWriter, r *http.Request) {
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

	// Fetch current question state
	var q model.Question
	err = QuestionsCollection.FindOne(ctx, filter).Decode(&q)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Question not found",
			StatusCode: http.StatusNotFound,
		})
		return
	}

	// Toggle `hidden` field
	q.Hidden = !q.Hidden
	update := bson.M{"$set": bson.M{"hidden": q.Hidden}}

	_, err = QuestionsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.Response{
			Message:    "Failed to update question",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	message := "Question is now visible"
	if q.Hidden {
		message = "Question is now hidden"
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{
		Message:    message,
		StatusCode: http.StatusOK,
		Data:       q,
	})
}
