package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

func startServer() {

	r := mux.NewRouter()

	http.Handle("/assets/", http.StripPrefix(
		"/assets/",
		http.FileServer(http.Dir("assets")),
	))

	r.HandleFunc("/", root).Methods("GET")
	r.HandleFunc("/data", data).Methods("GET")
	r.HandleFunc("/updateState", updateState).Methods("POST")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", r))
}

func root(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("./assets/html/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusFound)
	}
	if err := tmpl.ExecuteTemplate(w, tmpl.Name(), nil); err != nil {
		log.Fatalf("homeHandler: %+v", err)
	}
}

func data(w http.ResponseWriter, r *http.Request) {


}

func updateState(w http.ResponseWriter, r *http.Request) {

}
