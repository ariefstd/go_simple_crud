package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/mydb"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/read", readHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/delete", deleteHandler)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		age := r.FormValue("age")

		_, err := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", name, age)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/read", http.StatusSeeOther)
	} else {
		tmpl := template.Must(template.ParseFiles("templates/create.html"))
		tmpl.Execute(w, nil)
	}
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []struct {
		ID   int
		Name string
		Age  int
	}

	for rows.Next() {
		var user struct {
			ID   int
			Name string
			Age  int
		}
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	tmpl := template.Must(template.ParseFiles("templates/read.html"))
	tmpl.Execute(w, users)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		name := r.FormValue("name")
		age := r.FormValue("age")

		_, err := db.Exec("UPDATE users SET name = ?, age = ? WHERE id = ?", name, age, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/read", http.StatusSeeOther)
	} else {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Redirect(w, r, "/read", http.StatusSeeOther)
			return
		}

		var user struct {
			ID   int
			Name string
			Age  int
		}
		err := db.QueryRow("SELECT id, name, age FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/update.html"))
		tmpl.Execute(w, user)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")

		_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/read", http.StatusSeeOther)
	} else {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Redirect(w, r, "/read", http.StatusSeeOther)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/delete.html"))
		tmpl.Execute(w, map[string]string{"ID": id})
	}
}
