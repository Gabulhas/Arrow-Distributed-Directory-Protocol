package Nodes

import "fmt"

type Node struct {
	Type       NodeType //Tipo do Node, ver NodeType.go
	MyChan     string   //Channel onde recebe acesso ao objeto
	Find       string   //Channel onde recebe pedidos
	Link       string   //Ligação para o child Node
	WaiterChan string   //Channel de quem fez pedido
	MyAddress  string   //Endereço do Node
	VisAddress string   //Endereço onde faz "update" do seu estado atual para visualização
	Obj        bool 	//Se tem objeto ou não (redundante)
}

func (node *Node) OutputState() {
	fmt.Printf("\n---------------------State-------------------")
	fmt.Printf("\nMy Address:%s", node.MyAddress)
	fmt.Printf("\nLink:%s", node.Link)
	fmt.Printf("\nWaiter Chan:%s", node.WaiterChan)
	fmt.Printf("\nCurrent State : %s", node.Type.String())
	fmt.Printf("\n---------------------------------------------")
}
