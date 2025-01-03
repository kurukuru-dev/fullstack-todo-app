package main

import (
	"backend/app/database"
	"backend/app/handler"
	"backend/app/middleware"
	"backend/app/router"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("failed to initialize DB: %v", err)
	}
	defer database.GetDB().Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/todos", router.MethodRouter(map[string]http.HandlerFunc{
		http.MethodGet:  handler.GetTodos,
		http.MethodPost: handler.CreateTodo,
	}))
	mux.HandleFunc("/todos/", router.MethodRouter(map[string]http.HandlerFunc{
		http.MethodGet:    handler.GetTodoById,
		http.MethodPut:    handler.UpdateTodoById,
		http.MethodDelete: handler.DeleteTodoById,
	}))

	lateLimiter := middleware.NewRateLimiter()
	handlerWithMiddlewares := middleware.Chain(mux, lateLimiter)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithMiddlewares))
}
