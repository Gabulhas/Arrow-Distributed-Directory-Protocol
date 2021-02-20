package main

import (
	"flag"
	"fmt"
	"os"
	"projeto/Controller"
	"projeto/Nodes"
	"projeto/utils"
)

var selfNode *Nodes.Node

func init() {

	selfNode = new(Nodes.Node)

	myAddress := flag.String("address", "", "Node's Address (Required)")
	myType := flag.Int("type", -1, "Node's Type (0-4)(Required)") //Por definição de todos os tipos
	visAddress := flag.String("visualization", "", "Visualization address.")
	link := flag.String("link", "", "Link")
	requests := flag.Bool("requests", true, "If this Node, when Idle, preforms Object Requests")

	flag.Parse()

	//TODO: se for IDLE é necessário o LINK!
	if *myType < 0 || *myType > 4 || *myAddress == ""{
		flag.PrintDefaults()
		os.Exit(-1)
	}



	selfNode.MyAddress = *myAddress
	selfNode.MyChan = fmt.Sprintf("http://%s/myChan", *myAddress)
	selfNode.Find = fmt.Sprintf("http://%s/find", *myAddress)

	selfNode.Type = Nodes.NodeType(*myType)



	selfNode.VisAddress = fmt.Sprintf("http://%s", *visAddress)

	selfNode.Link = *link
	if *link != "" {
		selfNode.Link = fmt.Sprintf("http://%s/find", *link)
	}

	if selfNode.Type == Nodes.OWNER_TERMINAL || selfNode.Type == Nodes.OWNER_WITH_REQUEST {
		selfNode.Obj = true
	} else if selfNode.Type == Nodes.IDLE && *requests{
		go selfNode.AutoRequest()
	}

	fmt.Println(utils.StructToString(selfNode))
}

func main() {
	selfNode.OutputState()

	go selfNode.UpdateVisualization()
	go ShellStart()

	Controller.StartServer(selfNode)
}
