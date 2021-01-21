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

var re = regexp.MustCompile(`http://|/find|/myChan`)
var Mutex sync.Mutex
var requestHistory []string
var queueHistory []string
var ownerHistory []string
var currentOwner = ""

func startServer() {

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	r.HandleFunc("/", root).Methods("GET")
	r.HandleFunc("/data", data).Methods("GET")
	r.HandleFunc("/queue", queue).Methods("GET")
	r.HandleFunc("/updateState", updateState).Methods("POST")
	r.HandleFunc("/logs", getLogs).Methods("GET")
	r.HandleFunc("/requestAll", requestAll).Methods("GET")

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
	previous := Nodes[update.MyAddress]
	Nodes[update.MyAddress] = update

	//Fez request e mudou para waiter
	if update.Type == 4 {
		requestHistory = append(requestHistory, update.MyAddress)
	}

	//Era waiter e recebeu acesso ao objeto
	if (previous.Type == 3 && update.Type == 0) || (previous.Type == 4 && update.Type == 1) && currentOwner != update.MyAddress {
		currentOwner = update.MyAddress
		if len(ownerHistory) == 0 {
			ownerHistory = append(ownerHistory, currentOwner)
		} else if ownerHistory[len(ownerHistory)-1] != currentOwner {
			ownerHistory = append(ownerHistory, currentOwner)
		}
	}

	Mutex.Unlock()

	json.NewEncoder(w).Encode("Successful")

}

func queue(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	response := new(elements.QueueResponse)

	var currentNode elements.Node
	var nextNode elements.Node

	Mutex.Lock()
	response.Requesting = requestHistory
	response.OwnerHistory = ownerHistory
	response.CurrentOwner = currentOwner
	currentNode = Nodes[currentOwner]

	if currentOwner == "" {
		Mutex.Unlock()
		return
	}

	for currentNode.WaiterChan != "" {
		nextNode = Nodes[re.ReplaceAllString(currentNode.WaiterChan, ``)]
		response.QueueNodes = append(response.QueueNodes, nextNode.MyAddress)
		currentNode = nextNode
	}

	//dever haver algoritmo mais simples que este
	pivot := 0
	flag := false
	if len(queueHistory) != 0 {

		for i := 0; i < len(queueHistory); i++ {
			if pivot < len(response.QueueNodes) {
				if queueHistory[i] == response.QueueNodes[pivot] {
					pivot = pivot + 1
					flag = true
				} else if flag {
					break
				}

			}
		}

		for i := pivot; i < len(response.QueueNodes); i++ {
			response.QueueHistory = append(response.QueueHistory, response.QueueNodes[i])
		}
	}

	requestHistory = nil
	ownerHistory = nil
	queueHistory = response.QueueNodes
	Mutex.Unlock()

	json.NewEncoder(w).Encode(response)
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

func requestAll(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var wg sync.WaitGroup

	w.Header().Set("Content-Type", "text/plain")

	for _, element := range Nodes {
		wg.Add(1)
		go remoteRequest(element.MyAddress, &wg)
	}
	wg.Wait()
	w.Write([]byte("Successful"))
}

// TODO: usar esta função para o remote request de um só node, do "double click"
func remoteRequest(address string, wg *sync.WaitGroup) {

	defer (*wg).Done()

	_, err := http.Get(fmt.Sprintf("http://%s/remoteRequest", address))
	if err != nil {
		fmt.Println(err)
	}
}
