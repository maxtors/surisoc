package surisoc

import "encoding/json"

// SocketResponse is a json message recived from the suricata socket
type SocketResponse struct {
	Return  string      `json:"return"`
	Message interface{} `json:"message"`
}

// ToString converts the response message element to a string
func (sr *SocketResponse) ToString() (string, error) {
	bytes, err := json.MarshalIndent(sr.Message, "", "    ")
	if err != nil {
		return "", nil
	}
	return string(bytes), nil
}
