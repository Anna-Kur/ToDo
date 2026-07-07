package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"Sprint-13-14/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "Title is required", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, "Failed to add task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id": id,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(dateFormat)

	if task.Date == "" {
		task.Date = today
		return nil
	}

	if task.Date == "today" || task.Date == today {
		task.Date = today
		return nil
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	if t.After(now) {
		return nil
	}

	if task.Repeat != "" {
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

		task.Date = nextDate

	} else {
		task.Date = today
	}

	return nil

}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "ID is required", http.StatusBadRequest)
		return
	}

	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeJSONError(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if err.Error() == "Task not found" {
			writeJSONError(w, "Task not found", http.StatusNotFound)
		} else {
			writeJSONError(w, "Failed to get task: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeJSONError(w, "ID is required", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "Title is required", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		if err.Error() == "Task not found" {
			writeJSONError(w, "Task not found", http.StatusNotFound)
		} else {
			writeJSONError(w, "Failed to update task: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "ID is required", http.StatusBadRequest)
		return
	}

	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeJSONError(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		if err.Error() == "Task not found" {
			writeJSONError(w, "Task not found", http.StatusNotFound)
		} else {
			writeJSONError(w, "Failed to delete task: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	response := map[string]string{
		"error": errorMsg,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
