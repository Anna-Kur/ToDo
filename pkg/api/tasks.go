package api

import (
	"Sprint-13-14/pkg/db"
	"encoding/json"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSONError(w, "Failed to get tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}

	response := TasksResp{
		Tasks: tasks,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
