package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/AniketGodambe/mongoapi/model"
)

var (
	questions  = make(map[int]model.Question)
	questionID = 1
	mutex      sync.Mutex
)

func AddQuestion(w http.ResponseWriter, r *http.Request) {
	var q model.Question
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Invalid request", nil)
		return
	}
	mutex.Lock()
	q.ID = questionID
	questionID++
	questions[q.ID] = q
	mutex.Unlock()
	respondWithJSON(w, http.StatusCreated, "Question added successfully", q)
}

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Invalid question ID", nil)
		return
	}
	var q model.Question
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Invalid request", nil)
		return
	}
	mutex.Lock()
	if _, exists := questions[id]; !exists {
		mutex.Unlock()
		respondWithJSON(w, http.StatusNotFound, "Question not found", nil)
		return
	}
	q.ID = id
	questions[id] = q
	mutex.Unlock()
	respondWithJSON(w, http.StatusOK, "Question updated successfully", q)
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Invalid question ID", nil)
		return
	}
	mutex.Lock()
	if _, exists := questions[id]; !exists {
		mutex.Unlock()
		respondWithJSON(w, http.StatusNotFound, "Question not found", nil)
		return
	}
	delete(questions, id)
	mutex.Unlock()
	respondWithJSON(w, http.StatusOK, "Question deleted successfully", nil)
}

func HideQuestion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Invalid question ID", nil)
		return
	}
	mutex.Lock()
	if q, exists := questions[id]; exists {
		q.Hidden = !q.Hidden // Toggle the hidden state
		questions[id] = q
		mutex.Unlock()
		message := "Question is now visible"
		if q.Hidden {
			message = "Question is now hidden"
		}
		respondWithJSON(w, http.StatusOK, message, q)
	} else {
		mutex.Unlock()
		respondWithJSON(w, http.StatusNotFound, "Question not found", nil)
	}
}

func respondWithJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.Response{Message: message, StatusCode: status, Data: data})
}
