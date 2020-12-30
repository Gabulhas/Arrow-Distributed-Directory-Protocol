package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"visualization/elements"
)

//TODO: Mudar para package

var re = regexp.MustCompile(`http://|/find`)
var Mutex sync.Mutex

func startServer() {

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	r.HandleFunc("/", root).Methods("GET")
	r.HandleFunc("/data", data).Methods("GET")
	r.HandleFunc("/updateState", updateState).Methods("POST")
	r.HandleFunc("/logs", getLogs).Methods("GET")

	log.Fatal(http.ListenAndServe(os.Getenv("address"), r))
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
	defer r.Body.Close()

	response := new(elements.VisResponse)

	var tempNodes []elements.Node
	var tempLinks []elements.Link

	for _, v := range Nodes {
		tempNodes = append(tempNodes, v)

		if v.Link != "" {
			tempLinks = append(tempLinks, elements.Link{
				Source: v.MyAddress,
				Target: v.Link,
			})
		}
	}

	response.Nodes = tempNodes
	response.Links = tempLinks

	json.NewEncoder(w).Encode(response)
}

func updateState(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var update elements.Node
	_ = json.NewDecoder(r.Body).Decode(&update)

	Mutex.Lock()
	AllUpdates = append(AllUpdates, update)
	update.Link = re.ReplaceAllString(update.Link, ``)
	Nodes[update.MyAddress] = update
	Mutex.Unlock()

	json.NewEncoder(w).Encode("Successful")

}

func getLogs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "text/plain")
	var sb strings.Builder
	sb.Grow(len(AllUpdates))
	for i, elem := range AllUpdates {
		fmt.Fprintf(&sb, "\n%d: %s, %s, %d", i, elem.MyAddress, elem.Link, elem.Type)
	}
	w.Write([]byte(sb.String()))

}
