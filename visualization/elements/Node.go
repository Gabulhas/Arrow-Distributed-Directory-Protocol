package elements

type Node struct {
	Type       int    `json:"Type"`
	MyAddress  string `json:"MyAddress"`
	Link       string `json:"Link"`
	WaiterChan string `json:"WaiterChan"`
}
