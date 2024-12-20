package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// DBの初期化
func initDB() error {
    var err error
    // dsn -> ユーザー名:パスワード@tcp(ホスト名:ポート番号)/データベース名?オプション
    dsn := "todo_user:todo_password@tcp(mysql-container:3306)/todo_db"
    db, err = sql.Open("mysql", dsn) // DBとの接続を準備
    if err != nil {
        return fmt.Errorf("failed to connect to DB: %w", err)
    }

    // DBへの接続を確認
    if err := db.Ping(); err != nil {
        return fmt.Errorf("failed to ping DB: %w", err)
    }

    log.Println("Connected to the database successfully!")
	return nil
}

type Todo struct {
    Id int `json:"id"`
    Title string `json:"title"`
    IsComplete bool `json:"is_completed"`
}

// Todoリストをすべて取得する
func getTodos(w http.ResponseWriter, _ *http.Request) {
    rows, err := db.Query("SELECT id, title, is_completed FROM todos")
    if err != nil {
        http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
        return
    }

    var todos []Todo
    // レコードがある限り、次の行に進む
    for rows.Next() {
        var todo Todo
        if err := rows.Scan(&todo.Id, &todo.Title, &todo.IsComplete); err != nil {
            http.Error(w, "Failed to scan todos", http.StatusInternalServerError)
            return
        }
        todos = append(todos, todo)
    } 

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

// TodoリストのIDを指定して取得する
func getTodoById(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
    } 

    var todo Todo
    query := "SELECT id, title, is_completed FROM todos WHERE id = ?"
    err = db.QueryRow(query, id).Scan(&todo.Id, &todo.Title, &todo.IsComplete); 
    if err != nil {
        // QueryRow()は結果がない場合、sql.ErrNoRowsを返すため、適切なエラーハンドリングを行う
        if err == sql.ErrNoRows {
            http.Error(w, "Todo not found", http.StatusNotFound)
        } else {
            http.Error(w, "Failed to fetch todo", http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}

// Todoリストを追加する
func createTodo(w http.ResponseWriter, r *http.Request) {
    var newTodo Todo
    if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    results, err := db.Exec("INSERT INTO todos (title, is_completed) VALUES (?, ?)", newTodo.Title, newTodo.IsComplete)
    if err != nil {
        http.Error(w, "Failed to insert todo", http.StatusInternalServerError)
        return
    }
    
    id, err := results.LastInsertId()
    if err != nil {
        http.Error(w, "Failed to get last insert ID", http.StatusInternalServerError)
        return
    }

    // LastInsertId()はint64型を返すので、int型に変換
    newTodo.Id = int(id)
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

    query := "UPDATE todos SET title = ?, is_completed = ? WHERE id = ?"
    result, err := db.Exec(query, updatedTodo.Title, updatedTodo.IsComplete, id)
    if err != nil {
        http.Error(w, "Failed to update todo", http.StatusInternalServerError)
        return
    }

    // ResultインターフェースのRowsAffected()は更新された行数を返す。
    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        http.Error(w, "Todo not found", http.StatusNotFound)
        return
    }

    updatedTodo.Id = id
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedTodo)
}

// TodoリストのIDを指定して削除する
func deleteTodoById(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    result, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
    if err != nil {
        http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        http.Error(w, "Todo not found", http.StatusNotFound)
        return
    }

    // 削除が成功した場合、ステータスコード204を返す
    w.WriteHeader(http.StatusNoContent)
}

func main() {
    if err := initDB(); err != nil {
        log.Fatalf("failed to initialize DB: %v", err)
    }
    defer db.Close()

    http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTodos(w, r) 
		case http.MethodPost:
			createTodo(w, r) 
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

    http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTodoById(w, r)
		case http.MethodPut:
			updateTodoById(w, r) 
		case http.MethodDelete:
			deleteTodoById(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

    log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}