package elements

type QueueResponse struct {
	QueueNode  []string `json:"queue_nodes"`
	Owner      string   `json:"owner"`
	Requesting []string `json:"requesting"`
}
