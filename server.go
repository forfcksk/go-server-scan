package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Socks struct {
	ID   int    `json:"id"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:mysrdbz2KW@tcp(192.168.56.1:3306)/socks")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/socks", getPosts).Methods("GET")
	router.HandleFunc("/socks", createPost).Methods("POST")
	router.HandleFunc("/socks/{id}", getPost).Methods("GET")
	router.HandleFunc("/socks/{id}", deletePost).Methods("DELETE")
	router.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":8000", router)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var posts []Socks

	result, err := db.Query("SELECT id, ip, port FROM socks")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var post Socks
		err := result.Scan(&post.ID, &post.IP, &post.Port)
		if err != nil {
			panic(err.Error())
		}
		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stmt, err := db.Prepare("INSERT INTO socks (ID, IP, Port) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newid := keyVal["ID"]
	newip := keyVal["IP"]
	newport := keyVal["Port"]

	_, err = stmt.Exec(newid, newip, newport)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New post was created")
}

func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT ip, port FROM socks WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var post Socks
	for result.Next() {
		err := result.Scan(&post.IP, &post.Port)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(post)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM socks WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Post with ID = %s was deleted", params["id"])
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	result, err := db.Query("SELECT id, ip, port FROM socks")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()
	proxyes := []Socks{}

	for result.Next() {
		s := Socks{}
		err := result.Scan(&s.ID, &s.IP, &s.Port)
		if err != nil {
			panic(err.Error())
		}
		proxyes = append(proxyes, s)
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, proxyes)
}
