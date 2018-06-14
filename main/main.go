package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)


type Person struct {
    ID        string   `json:"id,omitempty"`
    Firstname string   `json:"firstname,omitempty"`
    Lastname  string   `json:"lastname,omitempty"`
   // Address   *Address `json:"address,omitempty"`
}
/* type Address struct {
    City  string `json:"city"`
    State string `json:"state"`
} */

//var people []Person

func main(){	
	router := mux.NewRouter()
	router.HandleFunc("/", Default).Methods("GET")
	router.HandleFunc("/people", GetPeople).Methods("GET")
    router.HandleFunc("/person/{id}", GetPerson).Methods("GET")
    router.HandleFunc("/person/", CreatePerson).Methods("POST")
    router.HandleFunc("/person/{id}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))

	
	
}
//Default page
func Default(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("Hello World!"))
}

//Get all people
 func GetPeople(w http.ResponseWriter, r *http.Request) {
	database, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(500)
		return
	}
	rows, err := database.Query("SELECT id, firstname, lastname FROM people")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to create table"))
		return
	}
	defer rows.Close()
	var people []Person
	for rows.Next() {
		var person Person
		if err := rows.Scan(
			&person.ID, &person.Firstname,&person.Lastname,
		); err != nil {
			log.Fatal(err)
			return
		}
		people = append(people, person)
	}
	json.NewEncoder(w).Encode(people)
} 


//Get person with specific id
 func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(500)
		return
	}
	defer db.Close()
	rows, err := db.Query(`select id, firstname, lastname from people where id = ?;`, params["id"])

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed!"))
		return
	}

	
	var person Person
	for rows.Next() {
		
		if err := rows.Scan(
			&person.ID, &person.Firstname,&person.Lastname,
		); err != nil {
			log.Fatal(err)
			return
		}
	}
	json.NewEncoder(w).Encode(person) 
} 

func CreatePerson(w http.ResponseWriter, r *http.Request) {
    var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	log.Println("hello")
	log.Println(person.Firstname)
	database, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(500)
		return
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to create table"))
		return
	}
	statement.Exec()
	statement, err = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed!"))
		return
	}
	statement.Exec(person.Firstname, person.Lastname)
	json.NewEncoder(w).Encode(person)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	//DELETE FROM employees
//WHERE last_name = 'Smith';
	params := mux.Vars(r)
    /* for index, item := range people {
        if item.ID == params["id"] {
			log.Println("item to delete")
			people = append(people[:index], people[index+1:]...)
            		break
		}
	}
	json.NewEncoder(w).Encode(people) */



	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(500)
		return
	}
	defer db.Close()
	rows, err := db.Query(`delete from people where id = ?;`, params["id"])
	
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed!"))
		return
	}

	var person Person
	for rows.Next() {
		
		if err := rows.Scan(
			&person.ID, &person.Firstname,&person.Lastname,
		); err != nil {
			log.Fatal(err)
			return
		}
	}
	json.NewEncoder(w).Encode(person) 


} 
