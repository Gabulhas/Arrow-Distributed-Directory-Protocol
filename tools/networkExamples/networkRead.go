package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type NodeConfigs struct {
	MyAddress int
	NodeType  int
	Link      int
}

var indexRegex, _ = regexp.Compile(`<index>`)
var addressRegex, _ = regexp.Compile(`<address>`)
var typeRegex, _ = regexp.Compile(`<type>`)
var linkRegex, _ = regexp.Compile(`<link>`)
var vis_addressRegex, _ = regexp.Compile(`<vis_address>`)
var portRegex, _ = regexp.Compile(`<port>`)

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
      requests: "true"
    ports:
      - "<port>:<port>"
    network_mode: host


`

var nodes []NodeConfigs
var vis_address string

/*


Turns CSV files into docker-compose.yml files based on the templates above
Args: 1 - name of .csv file inside the csv_files directory
	  2 - visualization address
*/

func main() {

	vis_address = os.Args[2]

	readCSV()
	toDockerCompose()
	writeToFile()
}

func readCSV() {
	csvFile, err := os.Open(fmt.Sprintf("./csv_files/%s", os.Args[1]))
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
	//TODO:Mudar para string builder
	for i, node := range nodes {
		service := templateDockerCompose
		port := 8000 + node.MyAddress
		linkValue := 8000 + node.Link
		link := fmt.Sprintf("127.0.0.1:%d", linkValue)

		service = indexRegex.ReplaceAllString(service, strconv.Itoa(i))
		service = addressRegex.ReplaceAllString(service, fmt.Sprintf("127.0.0.1:%d", port))
		service = typeRegex.ReplaceAllString(service, fmt.Sprintf("%d", node.NodeType))

		if node.Link == -1 {
			link = ""
		}
		service = linkRegex.ReplaceAllString(service, link)
		service = vis_addressRegex.ReplaceAllString(service, vis_address)
		service = portRegex.ReplaceAllString(service, fmt.Sprintf("%d", port))

		DockerComposeStart = DockerComposeStart + service
	}

}

func writeToFile() {

	f, err := os.Create(fmt.Sprintf("./dockerfiles/%s", strings.Replace(os.Args[1], ".csv", ".yml", 1)))
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
