package Nodes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"projeto/Channels"
	"time"
)

var max_retries = 4

func (node *Node) SendThroughLink(accessRequest Channels.AccessRequest) {
	go sendDataTo(node.Link, accessRequest)
}

func (node *Node) SendObjectAccess(giveAccess Channels.GiveAccess) {
	go sendDataTo(node.WaiterChan, giveAccess)
}

func sendDataTo(channel string, data interface{}) {

	message, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	retries := 0

	for retries < max_retries {

		resp, err := http.Post(channel, "application/json", bytes.NewBuffer(message))
		if err != nil {
			fmt.Printf(err.Error() + "")

			retries++
			time.Sleep(time.Second * time.Duration(3))

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

func (node *Node) UpdateVisualization() {

	message, err := json.Marshal(node)
	if err != nil {
		log.Fatal(err)
	}

	retries := 0

	for retries < max_retries {
		resp, err := http.Post(node.VisAddress, "application/json", bytes.NewBuffer(message))
		if err != nil {
			fmt.Printf(err.Error() + "")
			retries++
			time.Sleep(time.Second * time.Duration(3))
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
