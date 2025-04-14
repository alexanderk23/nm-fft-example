package client

import "encoding/json"

type PortValues struct {
	Port   uint16
	Values []uint16
}

func (r *PortValues) MarshalJSON() ([]byte, error) {
	arr := []any{r.Port, r.Values}
	return json.Marshal(arr)
}

type PortsMessage struct {
	Ports []PortValues `json:"ports"`
}
