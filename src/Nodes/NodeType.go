package Nodes

type NodeType int

const (
	OWNER_WITH_REQUEST NodeType = iota
	OWNER_TERMINAL
	IDLE
	WAITER_WITH_REQUEST
	WAITER_TERMINAL
)

// String() devolve Nome do Tipo
func (nodeType NodeType) String() string {
	return [...]string{"Owner With Request", "Owner Terminal", "Idle", "Waiter With Request", "Waiter Terminal"}[nodeType]
}
