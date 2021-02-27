package elements

type VisResponse struct {
	Nodes     []Node       `json:"nodes"`
	Links     []Connection `json:"links"`
	QueueCons []Connection `json:"queue_cons"`
	//QueueResponse QueueResponse `json:"queue_response"`
}
