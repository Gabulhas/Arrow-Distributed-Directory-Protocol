package main

import (
	"fmt"
	"log"
	"os"
	"projeto/Controller"
	"projeto/Nodes"
	"strconv"
	"sync"
)

var node *Nodes.Node

var Mutex *sync.Mutex

func init() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Printf("Required arguments.")
		os.Exit(-1)
	}

	node = new(Nodes.Node)

	node.MyAddress = args[0]
	node.MyChan = fmt.Sprintf("http://%s/myChan", args[0])
	node.Find = fmt.Sprintf("http://%s/find", args[0])
	node.VisAddress = fmt.Sprintf("http://%s", os.Getenv("VIS_ADDRESS"))
	nodeType, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal("Failed parsing the 2nd argument")
	}
	node.Type = Nodes.NodeType(nodeType)

	if len(args) > 2 {
		node.Link = fmt.Sprintf("http://%s/find", args[2])
	}

	if node.Type == Nodes.OWNER_TERMINAL || node.Type == Nodes.OWNER_WITH_REQUEST {
		node.Obj = true
	} else if node.Type == Nodes.IDLE {
		go node.AutoRequest()
	}

}

func main() {
	node.OutputState()

	go node.UpdateVisualization()
	go ShellStart()

	Controller.StartServer(node)
}

