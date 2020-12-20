package main

import (
	"fmt"
	"projeto/Channels"
	"projeto/Nodes"
	"projeto/utils"
	"time"
)

func HandleFind(accessRequest Channels.AccessRequest) {
	Mutex.Lock()
	defer Mutex.Unlock()

	newAccessRequest := accessRequest

	switch Node.Type {
	case Nodes.OWNER_TERMINAL:
		Node.OwnerWithRequest(accessRequest.Link, accessRequest.GiveAccess.WaiterChan)
		go releaseObj()
		break
	case Nodes.OWNER_WITH_REQUEST:
		newAccessRequest.Link = Node.Find
		SendThroughLink(newAccessRequest)
		Node.OwnerWithRequest(accessRequest.Link, Node.WaiterChan)
		break
	case Nodes.IDLE:
		newAccessRequest.Link = Node.Find
		SendThroughLink(newAccessRequest)
		Node.Idle(accessRequest.Link)
		break
	case Nodes.WAITER_TERMINAL:
		Node.WaiterWithRequest(accessRequest.Link, accessRequest.GiveAccess.WaiterChan)
		break
	case Nodes.WAITER_WITH_REQUEST:

		newAccessRequest.Link = Node.Find
		SendThroughLink(newAccessRequest)

		Node.WaiterWithRequest(accessRequest.Link, Node.WaiterChan)
		break
	}
	fmt.Printf("new: %s old:%s", newAccessRequest.GiveAccess.WaiterChan, accessRequest.GiveAccess.WaiterChan)

	go utils.UpdateVisualization(*Node, visualization_address)

}

func releaseObj() {

	Mutex.Lock()

	randomSleep := utils.RandomRange(1, 2)

	fmt.Printf("\nReleasing the Object in %d seconds.", randomSleep)

	time.Sleep(time.Second * time.Duration(randomSleep))
	fmt.Printf("\nReleasing...")

	accessObject := Channels.GiveAccess{WaiterChan: Node.WaiterChan}
	SendObjectAccess(accessObject)
	Node.Idle(Node.Link)

	go autoRequest()
	Mutex.Unlock()

	go utils.UpdateVisualization(*Node, visualization_address)

}

func autoRequest() {
	randomSleep := utils.RandomRange(9, 20)

	fmt.Printf("\nRequesting the Object in %d seconds.", randomSleep)
	time.Sleep(time.Second * time.Duration(randomSleep))

	//decidir se faz pedido
	//chance de fazer
	if requests := utils.RandomRange(0, 3); requests > 0 {
		Request()
	} else {
		cooldown := utils.RandomRange(3, 8)
		time.Sleep(time.Second * time.Duration(cooldown))
		autoRequest()
	}

}

func Request() {
	Mutex.Lock()
	defer Mutex.Unlock()

	if Node.Type != Nodes.IDLE {
		fmt.Printf("Can't request an object if not Idle.")
		return
	}

	accessRequest := Channels.AccessRequest{
		GiveAccess: Channels.GiveAccess{
			WaiterChan: Node.MyChan,
		},
		Link: Node.Find,
	}
	SendThroughLink(accessRequest)
	Node.WaiterTerminal()

	go utils.UpdateVisualization(*Node, visualization_address)
}

//Adicionar node externo que contém o objeto
func ReceiveObj(giveAccess Channels.GiveAccess) {
	Mutex.Lock()
	defer Mutex.Unlock()

	switch Node.Type {
	case Nodes.WAITER_TERMINAL:
		Node.OwnerTerminal()
		break
	case Nodes.WAITER_WITH_REQUEST:
		Node.OwnerWithRequest(Node.Link, Node.WaiterChan)
		go releaseObj()
		break
	}

	go utils.UpdateVisualization(*Node, visualization_address)

}
