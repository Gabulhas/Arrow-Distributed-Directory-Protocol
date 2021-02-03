package Nodes

func (node *Node) OwnerWithRequest(newLink string, waiterChan string) {
	node.Type = OWNER_WITH_REQUEST
	node.Link = newLink
	node.Obj = true
	node.WaiterChan = waiterChan
}

func (node *Node) OwnerTerminal() {
	node.Type = OWNER_TERMINAL
	node.Link = ""
	node.Obj = true
	node.WaiterChan = ""

}

func (node *Node) Idle(newLink string) {
	node.Type = IDLE
	node.Link = newLink
	node.Obj = false
	node.WaiterChan = ""
}

func (node *Node) WaiterWithRequest(newLink string, waiterChan string) {
	node.Type = WAITER_WITH_REQUEST
	node.Link = newLink
	node.WaiterChan = waiterChan
}

func (node *Node) WaiterTerminal() {
	node.Type = WAITER_TERMINAL
	node.Link = ""
	node.WaiterChan = ""
}
