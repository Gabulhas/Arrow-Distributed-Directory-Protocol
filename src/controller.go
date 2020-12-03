package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"projeto/Channels"
	"projeto/utils"
)

func startServer() {
	r := mux.NewRouter()
	r.HandleFunc("/find", findRoute).Methods("POST")
	r.HandleFunc("/mychan", myChanRoute).Methods("POST")
	http.ListenAndServe(Node.MyAddress, r)
}
func findRoute(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var accessRequest Channels.AccessRequest
	_ = json.NewDecoder(r.Body).Decode(&accessRequest)
	fmt.Printf("\nGot a find request")

	find <- accessRequest

	json.NewEncoder(w).Encode("Successful")

}

func myChanRoute(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var giveAccess Channels.GiveAccess
	_ = json.NewDecoder(r.Body).Decode(&giveAccess)
	fmt.Printf("\nGot Access To The Object!")

	myChan <- giveAccess

	json.NewEncoder(w).Encode("Successful")
}

func SendThroughLink(accessRequest Channels.AccessRequest) {

	fmt.Printf("Sending %s", utils.StructToString(accessRequest))

	message, err := json.Marshal(accessRequest)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(Node.Link, "application/json", bytes.NewBuffer(message))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(string(body))
}

//Isto poderá ser simplificado, pois estas duas funções têm o mesmo corpo, usar interface{}
func SendObjectAccess(giveAccess Channels.GiveAccess) {

	fmt.Printf("Sending %s", utils.StructToString(giveAccess))

	message, err := json.Marshal(giveAccess)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(Node.WaiterChan, "application/json", bytes.NewBuffer(message))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(string(body))
}

//esta parte também comunicará com a visualização
func OutputState() {
	fmt.Printf("Node: %s", utils.StructToString(Node))
	fmt.Printf("\nCurrent State : %s", Node.Type.String())
}
