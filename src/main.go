package main

import (
	"fmt"
	"log"
	"os"
	"projeto/Controller"
	"projeto/Nodes"
	"strconv"
)

var selfNode *Nodes.Node

func init() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Printf("Required arguments.")
		os.Exit(-1)
	}

	selfNode = new(Nodes.Node)

	selfNode.MyAddress = args[0]
	selfNode.MyChan = fmt.Sprintf("http://%s/myChan", args[0])
	selfNode.Find = fmt.Sprintf("http://%s/find", args[0])
	selfNode.VisAddress = fmt.Sprintf("http://%s", os.Getenv("VIS_ADDRESS"))
	nodeType, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal("Failed parsing the 2nd argument")
	}
	selfNode.Type = Nodes.NodeType(nodeType)

	if len(args) > 2 {
		selfNode.Link = fmt.Sprintf("http://%s/find", args[2])
	}

	if selfNode.Type == Nodes.OWNER_TERMINAL || selfNode.Type == Nodes.OWNER_WITH_REQUEST {
		selfNode.Obj = true
	} else if selfNode.Type == Nodes.IDLE {
		go selfNode.AutoRequest()
	}

}

func main() {
	selfNode.OutputState()

	go selfNode.UpdateVisualization()
	go ShellStart()

	Controller.StartServer(selfNode)
}

