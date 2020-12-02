package main

import (
	"fmt"
	"projeto/Channels"
	"projeto/Nodes"
	"projeto/utils"
)

func HandleFind(accessRequest Channels.AccessRequest) {
	newAccessRequest := accessRequest
	fmt.Printf("\nOld find: " + utils.Struct_to_string(accessRequest))

	switch Node.Type {
	case Nodes.OWNER_TERMINAL:
		Node.OwnerWithRequest(accessRequest.Link, accessRequest.GiveAccess.WaiterChan)
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

	fmt.Printf("\nNew find: " + utils.Struct_to_string(newAccessRequest))
}
