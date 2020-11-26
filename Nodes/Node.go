package Nodes

type NodeType int

const (
	OWNER_WITH_REQUEST NodeType = iota
	OWNER_TERMINAL
	IDLE
	WAITER_WITH_REQUEST
	WAITER_TERMINAL
)

type Node struct {
	Type       NodeType //Obvio
	MyChan     string   //Channel onde recebe acesso ao objeto
	Find       string   //Channel onde recebe pedidos
	Link       string   //Ligação para o child Node
	WaiterChan string   //Channel de quem fez pedido
	Obj        bool     //mudar mais tarde, só diz se tem ou não Obj
}

func (node *Node) OwnerWithRequest(newLink string, waiterChan string) {
	node.Type = OWNER_WITH_REQUEST
	node.Link = newLink
	node.WaiterChan = waiterChan
}

func (node *Node) Idle(newLink string) {
	node.Type = IDLE
	node.Link = newLink
}

func (node *Node) WaiterWithRequest(newLink string, waiterChan string){
	node.Type = WAITER_WITH_REQUEST
	node.Link = newLink
	node.WaiterChan = waiterChan
}