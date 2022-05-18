package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type Intern struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	About     string `json:"about"`
}

const (
	host     = "postgres_container"
	port     = 5432
	user     = "intern"
	password = "intern"
	dbname   = "intern"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user, password, host, port, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM intern")
	if err != nil {
		log.Fatal(err)
	}

	var people []Intern

	for rows.Next() {
		var intern Intern
		rows.Scan(&intern.FirstName, &intern.LastName, &intern.About)
		people = append(people, intern)
	}

	peopleBytes, _ := json.MarshalIndent(people, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var p Intern
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO intern (first_name, last_name, about) VALUES ($1, $2, $3)`
	_, err = db.Exec(sqlStatement, p.FirstName, p.LastName, p.About)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func main() {
	r := mux.NewRouter()
	fmt.Println("Hello Leroy Merlin")
	r.HandleFunc("/intern", POSTHandler).Methods("POST")
	r.HandleFunc("/intern", GETHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
