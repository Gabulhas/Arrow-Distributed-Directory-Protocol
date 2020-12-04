package main

import (
	"fmt"
	"projeto/Channels"
	"projeto/Nodes"
	"projeto/utils"
	"time"
)

func HandleFind(accessRequest Channels.AccessRequest) {
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
		autoRequest()
		break
	case Nodes.WAITER_TERMINAL:
		Node.WaiterWithRequest(accessRequest.Link, accessRequest.GiveAccess.WaiterChan)
		break
	case Nodes.WAITER_WITH_REQUEST:
		Node.WaiterWithRequest(accessRequest.Link, accessRequest.GiveAccess.WaiterChan)
		break
	}

}

func releaseObj() {

	randomSleep := utils.RandomRange(10, 40)

	fmt.Printf("\nReleasing the Object in %d seconds.", randomSleep)

	time.Sleep(time.Second * time.Duration(randomSleep))
	fmt.Printf("\nReleasing...")
	//Mudar aqui para incluir o node (que está fora do grafo) que contém o object (e iso passa a ser accesso a objeto)
	//Este waiter chan está errado, mudar

	accessObject := Channels.GiveAccess{WaiterChan: Node.WaiterChan}
	SendObjectAccess(accessObject)
	Node.Idle(Node.Link)
	autoRequest()

}

//Arranjar forma de correr isto sempre que o node passar a ser Idle, talvez com channels
func autoRequest() {
	go func() {
		utils.RandomSleep(time.Minute, 20, 100)
		Request()
	}()
}

//Isto também poderá funcionar manualmente, mas para simulação irá correr de forma aleatória, por agora
//Implementar forma de detetar mudanças de estado sem repetição de código
func Request() {
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
}

//Adicionar node externo que contém o objeto
func ReceiveObj(giveAccess Channels.GiveAccess) {

	switch Node.Type {
	case Nodes.WAITER_TERMINAL:
		Node.OwnerTerminal()
		break
	case Nodes.WAITER_WITH_REQUEST:
		Node.OwnerWithRequest(Node.Link, Node.WaiterChan)
		go releaseObj()
		break
	}

}
