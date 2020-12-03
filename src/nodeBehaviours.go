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
		go func() {
			utils.RandomSleep(time.Second, 0, 4)
			releaseObj()
		}()
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
		Node.WaiterWithRequest(accessRequest.Link, accessRequest.GiveAccess.WaiterChan)
		break
	}

}

func releaseObj() {
	fmt.Printf("Releasing the Object.")
	//Mudar aqui para incluir o node (que está fora do grafo) que contém o object (e iso passa a ser accesso a objeto)
	//Este waiter chan está errado, mudar
	accessObject := Channels.GiveAccess{WaiterChan: Node.WaiterChan}
	go SendObjectAccess(accessObject)
	Node.Idle(Node.Link)

}

//Arranjar forma de correr isto sempre que o node passar a ser Idle, talvez com channels
func autoRequest() {
	go func() {
		utils.RandomSleep(time.Second, 20, 100)
		Request()
	}()
}

//Isto também poderá funcionar manualmente, mas para simulação irá correr de forma aleatória
//Implementar forma de detetar mudanças de estado
func Request() {
	accessRequest := Channels.AccessRequest{
		GiveAccess: Channels.GiveAccess{
			WaiterChan: Node.MyChan,
		},
		Link: Node.Find,
	}
	go SendThroughLink(accessRequest)

}

func receiveObj() {

}
