package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func checkID(r *http.Request) (bool, string) {
	id := chi.URLParam(r, "id")
	_, ok := tasks[id]
	if !ok {
		return false, "Task not Found"
	}
	return true, id
}

func getPoint(w http.ResponseWriter, r *http.Request) {

	found, message := checkID(r)
	if !found {
		http.Error(w, message, http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	var taskMass []Task

	found, message := checkID(r)
	if !found {
		http.Error(w, message, http.StatusNoContent)
		return
	}

	for _, task := range tasks {
		taskMass = append(taskMass, task)
	}

	task, err := json.Marshal(taskMass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(task)

}

func postPoint(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func deletePoint(w http.ResponseWriter, r *http.Request) {
	found, message := checkID(r)
	if !found {
		http.Error(w, message, http.StatusNoContent)
		return
	}
	delete(tasks, message)
	w.WriteHeader(http.StatusOK)

}
func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getPoints)
	r.Post("/tasks", postPoint)
	r.Get("/tasks/{id}", getPoint)
	r.Delete("/tasks/{id}", deletePoint)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
