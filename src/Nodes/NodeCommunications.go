package Nodes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"projeto/Channels"
	"time"
)

var maxRetries = 4

func (node *Node) SendThroughLink(accessRequest Channels.AccessRequest) {
	go sendDataTo(node.Link, accessRequest)
}

func (node *Node) SendObjectAccess(giveAccess Channels.GiveAccess) {
	go sendDataTo(node.WaiterChan, giveAccess)
}

func (node *Node) UpdateVisualization(){
	go sendDataTo(node.VisAddress, node)
}

func sendDataTo(toURL string, data interface{}) {

	message, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	retries := 0

	for retries < maxRetries {

		resp, err := http.Post(toURL, "application/json", bytes.NewBuffer(message))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)

			retries++
			time.Sleep(time.Second * time.Duration(4))

			continue
		}

		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			log.Fatal(err)
		}

		resp.Body.Close()
		break

	}

}
