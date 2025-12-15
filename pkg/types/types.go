package types

import "time"

type Pod struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Image     string    `json:"image"`
	Status    string    `json:"status"`
	NodeName  string    `json:"nodeName"`
	CreatedAt time.Time `json:"createdAt"`
}

type Node struct {
	Name      string    `json:"name"`
	IPAddress string    `json:"ipAddress"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type PodList struct {
	Items []Pod `json:"items"`
}

type NodeList struct {
	Items []Node `json:"items"`
}
