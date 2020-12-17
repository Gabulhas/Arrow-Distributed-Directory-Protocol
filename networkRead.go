package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type NodeConfigs struct {
	MyAddress int
	NodeType  int
	Link      int
}

var template = "docker run --env address:<address> type=<type> link=<link> VIS_ADDRESS= <vis_address>/updateState -p <port>:<port> --interactive --tty"

var DockerComposeStart = `
version: "3.8"
services:
`

var templateDockerCompose = `

  node_<index>:
    tty: true
    stdin_open: true
    build:
      context: ./src
      dockerfile: Dockerfile
    environment:
      address: <address>
      type: <type> 
      link: <link>
      VIS_ADDRESS: <vis_address>/updateState
    ports:
      - "<port>:<port>"
    network_mode: host


`

var nodes []NodeConfigs
var vis_address string


/*
Args: 1 - .csv file
	  2 - visualization address
      3 - output file

 */


func main() {

	vis_address = os.Args[2]

	readCSV()
	toDockerCompose()
	writeToFile()
}

func readCSV() {
	csvFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		tempAddress, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatal(err)
		}

		tempType, err := strconv.Atoi(line[1])
		if err != nil {
			log.Fatal(err)
		}

		tempLink, err := strconv.Atoi(line[2])
		if err != nil {
			tempLink = -1
		}

		nodes = append(nodes, NodeConfigs{
			MyAddress: tempAddress,
			NodeType:  tempType,
			Link:      tempLink,
		})

	}
	fmt.Println(nodes)
}

func toDockerCompose() {
	for i, node := range nodes {
		service := templateDockerCompose
		port := 8000 + node.MyAddress
		linkValue := 8000 + node.Link
		link := fmt.Sprintf("127.0.0.1:%d", linkValue)

		//mudar isto para regex
		service = strings.Replace(service, "<index>", strconv.Itoa(i), 1)
		service = strings.Replace(service, "<address>", fmt.Sprintf("127.0.0.1:%d", port), 1)
		service = strings.Replace(service, "<type>", fmt.Sprintf("%d", node.NodeType), 1)

		if node.Link == - 1 {
			link = ""
		}
		service = strings.Replace(service, "<link>", link, 1)
		service = strings.Replace(service, "<vis_address>", vis_address, 1)
		service = strings.ReplaceAll(service, "<port>", fmt.Sprintf("%d", port))

		DockerComposeStart = DockerComposeStart + service
	}

}

func writeToFile() {

	f, err := os.Create(fmt.Sprintf("%s.yml", os.Args[3]))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(DockerComposeStart)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Done")
}
