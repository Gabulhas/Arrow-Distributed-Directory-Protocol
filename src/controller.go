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
	r.HandleFunc("/myChan", myChanRoute).Methods("POST")
	r.HandleFunc("/request", remoteRequest).Methods("GET")
	http.ListenAndServe(Node.MyAddress, r)
}
func findRoute(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var accessRequest Channels.AccessRequest
	_ = json.NewDecoder(r.Body).Decode(&accessRequest)
	fmt.Printf("\nGot a find request")
	fmt.Printf("\n%s", utils.StructToString(accessRequest))

	HandleFind(accessRequest)

	json.NewEncoder(w).Encode("Successful")

}

func myChanRoute(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var giveAccess Channels.GiveAccess
	_ = json.NewDecoder(r.Body).Decode(&giveAccess)
	fmt.Printf("\nGot Access To The Object!")
	fmt.Printf("\n%s", utils.StructToString(giveAccess))

	fmt.Printf("\nGot the object. %s", utils.StructToString(giveAccess))
	ReceiveObj(giveAccess)

	json.NewEncoder(w).Encode("Successful")
}

func SendThroughLink(accessRequest Channels.AccessRequest) {
	fmt.Printf("\nSending %s", utils.StructToString(accessRequest))
	go sendDataTo(Node.Link, accessRequest)

}

//Isto poderá ser simplificado, pois estas duas funções têm o mesmo corpo, usar interface{}
func SendObjectAccess(giveAccess Channels.GiveAccess) {
	go sendDataTo(Node.WaiterChan, giveAccess)
}

func sendDataTo(channel string, data interface{}) string {

	message, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(channel, "application/json", bytes.NewBuffer(message))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)

}

func remoteRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	Request()
	json.NewEncoder(w).Encode("Requested")
}
