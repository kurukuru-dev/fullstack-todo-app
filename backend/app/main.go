package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Todo struct {
    Id int `json:"id"`
    Title string `json:"title"`
    IsComplete bool `json:"is_completed"`
}

var todos = []Todo{
    {Id: 1, Title: "Learn Go", IsComplete: false},
    {Id: 2, Title: "Build a RESTful API", IsComplete: false},
}

// Todoリストをすべて取得する
func getTodos(w http.ResponseWriter, _ *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

// TodoリストのIDを指定して取得する
func getTodosById(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/todos/")

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
    } 

    for _, todo := range todos {
        if todo.Id == id {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(todo)
            return
        }
    }

    http.Error(w, "Todo not found", http.StatusNotFound)
}

// Todoリストを追加する
func createTodo(w http.ResponseWriter, r *http.Request) {
    var newTodo Todo
    if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    newTodo.Id = len(todos) + 1
    todos = append(todos, newTodo)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    json.NewEncoder(w).Encode(newTodo)
}

// TodoリストのIDを指定して更新する
func updateTodoById(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var updatedTodo Todo
    if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    for i, todo := range todos {
        if todo.Id == id {
            todos[i] = updatedTodo
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(updatedTodo)
            return
        }
    }
}

// TodoリストのIDを指定して削除する
func deleteTodoById(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    for i, todo := range todos {
        if todo.Id == id {
            todos = append(todos[:i], todos[i+1:]...)
            
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(todos)
            return
        }
    }

    http.Error(w, "Todo not found", http.StatusNotFound)
}

func main() {
    fmt.Println("Hello Go!")
}