package Channels

// #1
type AccessRequest struct {
	GiveAccess GiveAccess `json:"giveAccess"` // #2
	Link       string     `json:"link"`       // #1
}

// #2
type GiveAccess struct {
	WaiterChan string `json:"waiterChan"`
}
