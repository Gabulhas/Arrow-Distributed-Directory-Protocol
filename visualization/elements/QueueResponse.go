package elements

type QueueResponse struct {
	QueueNodes   []string `json:"queue_nodes"`
	OwnerHistory []string `json:"owner_history"`
	Requesting   []string `json:"requesting"`
	QueueHistory []string `json:"queue_history"`
	CurrentOwner string   `json:"current_owner"`
}
