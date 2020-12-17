package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"projeto/Nodes"
)

func UpdateVisualization(update Nodes.Node, vis_address string) string {

	message, err := json.Marshal(update)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(vis_address, "application/json", bytes.NewBuffer(message))
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
