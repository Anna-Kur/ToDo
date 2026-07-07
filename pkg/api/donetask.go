package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"Sprint-13-14/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSONError(w, "Failed to delete task: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}

	now := time.Now()
	nextDate, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSONError(w, "Failed to calculate next date: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateDate(nextDate, id); err != nil {
		writeJSONError(w, "Failed to update task date: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
