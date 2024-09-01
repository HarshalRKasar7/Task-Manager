package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Task struct {
	Id          int
	Name        string
	Description string
	Status      bool
	dueDate     string
	dueTime     string
}

var (
	tasks  []Task // list of task
	nextId int // counter of ids
)

// Create task 
func (t *Task) CreateTask(name string, description string, date string, duetime string) {
	t.Id = nextId + 1
	t.Name = name
	t.Description = description
	t.dueDate = date
	t.dueTime = duetime
	nextId++
}

// saveTasksToFile saves the tasks slice to a JSON file.
func saveTasksToFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(tasks)
	if err != nil {
		fmt.Printf("Error encoding tasks to JSON: %v\n", err)
	}
}

// loadTasksFromFile loads the tasks from a JSON file.
func loadTasksFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return // No tasks file yet, ignore the error
		}
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		fmt.Printf("Error decoding tasks from JSON: %v\n", err)
		return
	}

	// Update nextId to be one more than the highest existing task ID
	for _, task := range tasks {
		if task.Id >= nextId {
			nextId = task.Id + 1
		}
	}
}

// send get request to add task 
func taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		name := r.URL.Query().Get("name")
		description := r.URL.Query().Get("description")
		date := r.URL.Query().Get("date")
		time := r.URL.Query().Get("time")

		fmt.Printf("\nTask Info\nTask Name: %s\nDescription: %s\nDate: %s\nTime: %s\n", name, description, date, time)
		var task Task
		task.CreateTask(name, description, date, time)

		tasks = append(tasks, task)
		saveTasksToFile("data.json")

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// displayTasksHandler handles displaying all tasks
func displayTasksHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(tasks)
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

// removeTaskHandler handles removing a task via HTTP GET request
func removeTaskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        idStr := r.URL.Query().Get("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "Invalid task ID", http.StatusBadRequest)
            return
        }

        for i, task := range tasks {
            if task.Id == id {
                tasks = append(tasks[:i], tasks[i+1:]...)
                saveTasksToFile("data.json")
                http.Redirect(w, r, "/", http.StatusSeeOther)
                return
            }
        }

        http.Error(w, "Task not found", http.StatusNotFound)
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

// completeTaskHandler handles marking a task as completed via HTTP GET request
func completeTaskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        idStr := r.URL.Query().Get("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "Invalid task ID", http.StatusBadRequest)
            return
        }
		for i, task := range tasks {
            if task.Id == id {
                tasks[i].Status = true
                saveTasksToFile("data.json")
                http.Redirect(w, r, "/", http.StatusSeeOther)
                return
            }
        }

        http.Error(w, "Task not found", http.StatusNotFound)
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

func main() {
	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/", fs)
	http.HandleFunc("/task", taskHandler)
	http.HandleFunc("/remove", removeTaskHandler)
    http.HandleFunc("/complete", completeTaskHandler)
    http.HandleFunc("/display", displayTasksHandler)
	loadTasksFromFile("data.json")

	fmt.Println("Serving on http:://localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err.Error())
	}
}
