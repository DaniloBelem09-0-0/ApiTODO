package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Todo struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}


var db *gorm.DB

func initDB() {
	dsn := "root:123456@tcp(localhost:3306)/tododb?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}) 
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.AutoMigrate(&Todo{})
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todos []Todo
	db.Find(&todos)
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	db.Create(todo)
	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var todo Todo
	if err := db.First(&todo, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_ = json.NewDecoder(r.Body).Decode(&todo)
	db.Save(&todo)
	json.NewEncoder(w).Encode(todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var todo Todo
	if err := db.First(&todo, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	db.Delete(&todo)
	json.NewEncoder(w).Encode(todo)
}


func main() {

	initDB()
	r := mux.NewRouter()
	r.HandleFunc("/todos", getTodos).Methods("GET")
	r.HandleFunc("/todos", createTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")

	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", r)
}