package Controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"projeto/Channels"
	"projeto/Nodes"
	"projeto/utils"
)

var node *Nodes.Node

func StartServer(newNode *Nodes.Node) {
	node = newNode
	r := mux.NewRouter()
	r.HandleFunc("/find", findRoute).Methods("POST")
	r.HandleFunc("/myChan", myChanRoute).Methods("POST")
	if err := http.ListenAndServe(node.MyAddress, r); err != nil {
		log.Fatal(err)
	}
}
func findRoute(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var accessRequest Channels.AccessRequest
	_ = json.NewDecoder(r.Body).Decode(&accessRequest)
	fmt.Printf("\nGot a find request")
	fmt.Printf("\n%s", utils.StructToString(accessRequest))

	node.HandleFind(accessRequest)

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
	node.ReceiveObj(giveAccess)

	json.NewEncoder(w).Encode("Successful")
}

