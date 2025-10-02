package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks []Task

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Не удалось преобразовать в JSON", http.StatusInternalServerError)
		return
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Добро пожаловать в To-Do API!")
}

func loadTacks() error {
	file, err := os.Open("data/tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []Task{}
			return err
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		return err
	}
	return nil
}

func saveTasks() error {
	file, err := os.Create("data/tasks.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var newTask Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	if newTask.Title == "" {
		http.Error(w, "Поле 'title' обязательно", http.StatusBadRequest)
		return
	}

	nextID := 1
	for _, t := range tasks {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}

	newTask.ID = nextID

	tasks = append(tasks, newTask)

	json.NewEncoder(w).Encode(newTask)

	if err := saveTasks(); err != nil {
		log.Printf("Ошибка сохранения: %v", err)
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path

	if len(path) <= len("/tasks/") {
		http.Error(w, "Неверный путь", http.StatusBadRequest)
		return
	}

	idStr := path[len("/tasks/"):]

	var taskID int
	_, err := fmt.Sscanf(idStr, "%d", &taskID)
	if err != nil || taskID <= 0 {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	index := -1
	for i, t := range tasks {
		if t.ID == taskID {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	tasks = append(tasks[:index], tasks[index+1:]...)
	w.WriteHeader(http.StatusNoContent)

	if err := saveTasks(); err != nil {
		log.Printf("Ошибка сохранения: %v", err)
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	if r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path

	if len(path) <= len("/tasks/") {
		http.Error(w, "Неверный путь", http.StatusBadRequest)
		return
	}

	idStr := path[len("/tasks/"):]

	var taskID int
	_, err := fmt.Sscanf(idStr, "%d", &taskID)
	if err != nil || taskID <= 0 {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	index := -1
	for i, t := range tasks {
		if t.ID == taskID {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	var updatedData Task
	err = json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	if updatedData.Title != "" {
		tasks[index].Title = updatedData.Title
	}
	tasks[index].Done = updatedData.Done

	json.NewEncoder(w).Encode(tasks[index])

	if err := saveTasks(); err != nil {
		log.Printf("Ошибка сохранения: %v", err)
	}
}

func main() {

	if err := loadTacks(); err != nil {
		log.Fatal("Не удалось загрузить задачи:", err)
	}

	http.HandleFunc("/", home)

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSuffix(r.URL.Path, "/")

		if path == "/tasks" && r.Method == http.MethodGet {
			getTasks(w, r)
			return
		}
		if path == "/tasks" && r.Method == http.MethodPost {
			createTask(w, r)
			return
		}
		if len(path) > len("/tasks/") {
			if r.Method == http.MethodDelete {
				deleteTask(w, r)
				return
			}
			if r.Method == http.MethodPut {
				updateTask(w, r)
				return
			}
		}
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
