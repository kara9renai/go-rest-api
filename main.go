package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/kara9renai/go-rest-api/database"
	"github.com/kara9renai/go-rest-api/params"
)

func main() {
	const port = ":8080"
	const dPath = ".sqlite3/todo.db"

	Db, err := database.NewDB(dPath)
	if err != nil {
		log.Println(err)
	}
	defer Db.Close()

	http.HandleFunc("/todos", todoHandler)
	http.HandleFunc("/todos/", todoIdHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getAllTodos(w, r)
	case "POST":
		postNewTodo(w, r)
	}
}

func todoIdHandler(w http.ResponseWriter, r *http.Request) {
	params := params.GetRouteParams(r)

	if len(params) < 2 || len(params) > 3 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// itemのidをintで取得
	id, err := strconv.Atoi(params[1])
	if err != nil || id < 1 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if len(params) == 2 {
		updateTodo(id, w, r)
	} else if params[2] == "done" {
		deleteDoneTodos(w)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

type Item struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

// アイテムの全取得
func getAllTodos(w http.ResponseWriter, r *http.Request) {

	var items []Item
	rows, err := database.Db.Query("SELECT * FROM todos;")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	for rows.Next() {
		var item Item
		rows.Scan(&item.Id, &item.Name, &item.Done)
		items = append(items, item)
	}
	var buf bytes.Buffer
	e := json.NewEncoder(&buf)
	if err := e.Encode(items); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, buf.String())
}

// アイテムの新規追加
func postNewTodo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	d := json.NewDecoder(r.Body)
	err := d.Decode(&req)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	_, err = database.Db.Exec(`INSERT INTO todos (name, done) VALUES (?, ?)`, req.Name, false)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}

	w.WriteHeader(201)
}

func updateTodo(id int, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		deleteTodo(id, w)
	default:
		http.Error(w, http.StatusText(500), 500)
	}
}

// 指定されたidのtodoを削除
func deleteTodo(id int, w http.ResponseWriter) {
	_, err := database.Db.Exec(`DELETE FROM todos WHERE id = ?`, id)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
	w.WriteHeader(200)
}

// 実行済みのタスクを削除
func deleteDoneTodos(w http.ResponseWriter) {
	_, err := database.Db.Exec(`DELETE FROM todos WHERE done = true`)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
	w.WriteHeader(200)
}
