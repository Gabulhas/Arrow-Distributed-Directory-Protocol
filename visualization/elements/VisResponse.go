package elements

type VisResponse struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
	//QueueResponse QueueResponse `json:"queue_response"`
}
