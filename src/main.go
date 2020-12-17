package main

import (
	"fmt"
	"log"
	"os"
	"projeto/Nodes"
	"projeto/utils"
	"strconv"
	"sync"
)

var Node *Nodes.Node

var Mutex *sync.Mutex

var visualization_address string

func main() {
	Node = new(Nodes.Node)
	initNode()
	OutputState()

	Mutex = &sync.Mutex{}

	utils.UpdateVisualization(*Node, visualization_address)
	go ShellStart()
	startServer()
}

func initNode() {
	visualization_address = fmt.Sprintf("http://%s", os.Getenv("VIS_ADDRESS"))

	args := os.Args[1:]

	fmt.Printf("All args")
	fmt.Println(args)

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
	} else if Node.Type == Nodes.IDLE {
		go autoRequest()
	}
}

//esta parte também comunicará com a visualização
func OutputState() {
	fmt.Printf("\n---------------------State-------------------")
	fmt.Printf("\nMy Address:%s", Node.MyAddress)
	fmt.Printf("\nLink:%s", Node.Link)
	fmt.Printf("\nWaiter Chan:%s", Node.WaiterChan)
	fmt.Printf("\nCurrent State : %s", Node.Type.String())
	fmt.Printf("\n---------------------------------------------")
}
