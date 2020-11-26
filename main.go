package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"projeto/Channels"
	"projeto/Nodes"
)

var node Nodes.Node

func main() {



	r := mux.NewRouter()
	r.HandleFunc("/find", findHandler).Methods("POST")
	r.HandleFunc("/mycahn", myChanHandler).Methods("POST")
	http.ListenAndServe(":8081", r)

}

func findHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	defer outputState()

	w.Header().Set("Content-Type", "application/json")

	var accessRequest, newAccessRequest Channels.AccessRequest
	_ = json.NewDecoder(r.Body).Decode(&accessRequest)
	spew.Println(accessRequest)

	switch node.Type {
	case Nodes.OWNER_TERMINAL:
		node.OwnerWithRequest(accessRequest.Link, accessRequest.GiveAccess.WaiterChan)
		break
	case Nodes.OWNER_WITH_REQUEST:
		newAccessRequest.Link = node.Find
		newAccessRequest.GiveAccess = accessRequest.GiveAccess
		sendThroughLink(newAccessRequest)
		node.OwnerWithRequest(accessRequest.Link, node.WaiterChan)
		break
	case Nodes.IDLE:
		newAccessRequest.Link = node.Find
		newAccessRequest.GiveAccess = accessRequest.GiveAccess
		sendThroughLink(newAccessRequest)
		node.Idle(accessRequest.Link)
		break
	case Nodes.WAITER_TERMINAL:
		node.WaiterWithRequest(accessRequest.Link, accessRequest.GiveAccess.WaiterChan)
		break
	case Nodes.WAITER_WITH_REQUEST:
		node.WaiterWithRequest(accessRequest.Link, accessRequest.GiveAccess.WaiterChan)
		break
	}

}

func myChanHandler(w http.ResponseWriter, r *http.Request) {

	node.Obj = true
}

func sendThroughLink(accessRequest Channels.AccessRequest) {

	message, err := json.Marshal(accessRequest)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(node.Link, "application/json", bytes.NewBuffer(message))
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
func outputState() {
	nodeJSON, err := json.Marshal(node)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("node: %s", string(nodeJSON))

}
