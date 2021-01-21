package Nodes

import (
	"fmt"
	"math/rand"
	"projeto/Channels"
	"projeto/utils"
	"sync"
	"time"
)

var Mutex = &sync.Mutex{}

func init(){
	rand.Seed(time.Now().UnixNano())
}

func (node *Node) HandleFind(accessRequest Channels.AccessRequest) {
	Mutex.Lock()
	defer Mutex.Unlock()

	newAccessRequest := accessRequest

	switch node.Type {
	case OWNER_TERMINAL:
		node.OwnerWithRequest(accessRequest.Link, accessRequest.GiveAccess.WaiterChan)
		go node.releaseObj()
		break
	case OWNER_WITH_REQUEST:
		newAccessRequest.Link = node.Find
		node.SendThroughLink(newAccessRequest)
		node.OwnerWithRequest(accessRequest.Link, node.WaiterChan)
		break
	case IDLE:
		newAccessRequest.Link = node.Find
		node.SendThroughLink(newAccessRequest)
		node.Idle(accessRequest.Link)
		break
	case WAITER_TERMINAL:
		node.WaiterWithRequest(accessRequest.Link, accessRequest.GiveAccess.WaiterChan)
		break
	case WAITER_WITH_REQUEST:
		newAccessRequest.Link = node.Find
		node.SendThroughLink(newAccessRequest)
		node.WaiterWithRequest(accessRequest.Link, node.WaiterChan)
		break
	}

	go node.UpdateVisualization()

}

func (node *Node) releaseObj() {
	Mutex.Lock()
	defer Mutex.Unlock()

	randomSleep := utils.RandomRange(1, 2)

	fmt.Printf("\nReleasing the Object in %d seconds.", randomSleep)

	time.Sleep(time.Second * time.Duration(randomSleep))
	fmt.Printf("\nReleasing...")

	accessObject := Channels.GiveAccess{WaiterChan: node.WaiterChan}

	node.SendObjectAccess(accessObject)
	node.Idle(node.Link)

	go node.AutoRequest()
	go node.UpdateVisualization()

}

func (node *Node) AutoRequest() {
	var randomSleep int

	for {

		randomSleep = utils.RandomRange(5, 15)

		fmt.Printf("\nTrying to Request the Object in %d seconds.", randomSleep)
		time.Sleep(time.Second * time.Duration(randomSleep))

		//decidir se faz pedido
		//chance de fazer
		if requests := utils.RandomRange(0, 1); requests > 0 {
			node.Request()
			break
		} else {
			cooldown := utils.RandomRange(5, 20)
			fmt.Printf("Didn't request. Retrying in %d seconds.", cooldown)
			time.Sleep(time.Second * time.Duration(cooldown))
		}
	}

}

func (node *Node) Request() {
	Mutex.Lock()
	defer Mutex.Unlock()

	//Existe para evitar:
	//que ou o utilizador faça um request e o node já mudou de tipo
	//que se faça um request a partir do método do Node de pedidos remotos
	if node.Type != IDLE {
		fmt.Printf("Can't request an object if not Idle.")
		return
	}
	fmt.Printf("Requesting.")

	accessRequest := Channels.AccessRequest{
		GiveAccess: Channels.GiveAccess{
			WaiterChan: node.MyChan,
		},
		Link: node.Find,
	}

	node.SendThroughLink(accessRequest)
	node.WaiterTerminal()

	go node.UpdateVisualization()
}

func (node *Node) ReceiveObj(giveAccess Channels.GiveAccess) {
	Mutex.Lock()
	defer Mutex.Unlock()

	fmt.Printf("Received Access:")
	fmt.Println(giveAccess)
	switch node.Type {
	case WAITER_TERMINAL:
		node.OwnerTerminal()
		break
	case WAITER_WITH_REQUEST:
		node.OwnerWithRequest(node.Link, node.WaiterChan)
		go node.releaseObj()
		break
	}

	go node.UpdateVisualization()

}
