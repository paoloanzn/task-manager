package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"github.com/paoloanzn/task-manager/internal/tasks"
)

var (
	manager *tasks.Manager = tasks.NewManager()
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload Response) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func Health(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "server running",
		Data:    struct{}{},
	})
	return
}

func createTask(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	title := params.Get("title")

	if title == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	id, err := manager.NewTask(title)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: fmt.Sprintf("%v", err),
			Data:    struct{}{},
		})
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "task created",
		Data: map[string]uint32{
			"id": id,
		},
	})

	return
}

func getTask(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get("id")

	if id == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	u64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "id is not an valid integer",
			Data:    struct{}{},
		})
		return
	}
	u32 := uint32(u64)

	task, err := manager.GetTask(u32)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: fmt.Sprintf("%v", err),
			Data:    struct{}{},
		})
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "task retrieved",
		Data: map[string]interface{}{
			"task": task,
		},
	})

	return
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get("id")

	if id == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	u64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "id is not an valid integer",
			Data:    struct{}{},
		})
		return
	}
	u32 := uint32(u64)

	err = manager.DeleteTask(u32)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: fmt.Sprintf("%v", err),
			Data:    struct{}{},
		})
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "task deleted",
		Data:    struct{}{},
	})

	return
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get("id")
	title := params.Get("title")
	status := params.Get("status")

	if id == "" || title == "" || status == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	idU64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "id is not an valid integer",
			Data:    struct{}{},
		})
		return
	}
	idU32 := uint32(idU64)

	parsedStatus, err := strconv.Atoi(status)
	if (err != nil) || (!tasks.IsValidStatus(parsedStatus)) {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "status is not an valid integer",
			Data:    struct{}{},
		})
		return
	}

	err = manager.UpdateTask(idU32, title, tasks.TaskStatus(parsedStatus))
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: fmt.Sprintf("%v", err),
			Data:    struct{}{},
		})
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "task updated",
		Data:    struct{}{},
	})

	return
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", Health)
	r.HandleFunc("/task/create", createTask)
	r.HandleFunc("/task/get", getTask)
	r.HandleFunc("/task/delete", deleteTask)
	r.HandleFunc("/task/update", updateTask)
	log.Fatal(http.ListenAndServe(":8080", r))
}
