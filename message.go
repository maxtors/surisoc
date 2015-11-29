package surisoc

import (
	"fmt"
	"strconv"
)

// SocketMessage is a json message that will be sent to the suricata socket
type SocketMessage struct {
	Command   string                  `json:"command"`
	Arguments *map[string]interface{} `json:"arguments"`
}

// NewSocketMessage will take a command and some arguments and create a
// a new socket message struct
func NewSocketMessage(command string, arguments ...string) (*SocketMessage, error) {
	var err error
	message := SocketMessage{Command: command}
	if len(arguments) != 0 {
		err = message.parseArguments(command, arguments...)
		if err != nil {
			return nil, &Error{Message: fmt.Sprintf("Could not parse arguments: %s, %+v (%s)", command, arguments, err.Error())}
		}
	}
	return &message, nil
}

func (sm *SocketMessage) parseArguments(command string, arguments ...string) error {
	var err error
	argumentsMap := make(map[string]interface{})
	switch command {
	case "iface-stat":
		if len(arguments) == 1 {
			argumentsMap["iface"] = arguments[0]
			break
		}
		return &Error{Message: fmt.Sprintf("iface-stat should have one argument, found %d", len(arguments))}
	case "pcap-file":
		if len(arguments) > 1 {
			argumentsMap["filename"] = arguments[0]
			argumentsMap["output-dir"] = arguments[1]
			if len(arguments) == 3 {
				argumentsMap["tenant"], err = strconv.Atoi(arguments[2])
				if err != nil {
					return &Error{Message: fmt.Sprintf("Tenant ID is not an intager: %s", arguments[2])}
				}
			}
			break
		}
		return &Error{Message: fmt.Sprintf("pcap-file should have atleast 2 arguments, found %d", len(arguments))}
	case "conf-get":
		if len(arguments) == 1 {
			argumentsMap["variable"] = arguments[0]
			break
		}
		return &Error{Message: fmt.Sprintf("conf-get should have one argument, found %d", len(arguments))}
	case "unregister-tenant-handler":
		return nil
	case "register-tenant-handler":
		return nil
	case "unregister-tenant":
		return nil
	case "register-tenant":
		return nil
	default:
		return &Error{Message: fmt.Sprintf("The command \"%s\" with %d arguments is unkown", command, len(arguments))}
	}
	sm.Arguments = &argumentsMap
	return nil
}
