package main

import (
	"fmt"
	"log"
	"os"
	"projeto/Channels"
	"projeto/Nodes"
	"strconv"
)

var Node *Nodes.Node
var find chan Channels.AccessRequest
var myChan chan Channels.GiveAccess
var typeChange chan Nodes.NodeType

func main() {
	Node = new(Nodes.Node)
	initNode()
	OutputState()

	find = make(chan Channels.AccessRequest, 10)
	myChan = make(chan Channels.GiveAccess, 10)
	typeChange = make(chan Nodes.NodeType)

	go ChanHandler()
	go ShellStart()
	startServer()
}

func initNode() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Printf("Required arguments")
		os.Exit(-1)
	}

	Node.MyAddress = args[0]
	Node.MyChan = fmt.Sprintf("http://%s/myChan", args[0])
	Node.Find = fmt.Sprintf("http://%s/find", args[0])
	nodeType, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal("Failed parsing the 2nd argument")
	}
	Node.Type = Nodes.NodeType(nodeType)

	if len(args) > 2 {
		Node.Link = fmt.Sprintf("http://%s/find", args[2])
	}

	if Node.Type == Nodes.OWNER_TERMINAL || Node.Type == Nodes.OWNER_WITH_REQUEST {
		Node.Obj = true
	}
}
