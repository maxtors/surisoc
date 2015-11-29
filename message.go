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
func NewSocketMessage(command string) *SocketMessage {
	return &SocketMessage{Command: command}
}

// ParseArgumentsList will check if any arguments are to be parsed and add to the message
func (sm *SocketMessage) ParseArgumentsList(arguments ...string) error {
	if len(arguments) != 0 {
		err := sm.parseArgumentsList(arguments...)
		if err != nil {
			return &Error{Message: fmt.Sprintf("Could not parse arguments: %s, %+v (%s)", sm.Command, arguments, err.Error())}
		}
	}
	return nil
}

// ParseArgumentsMap will parse a map of arguments and add to the message
func (sm *SocketMessage) ParseArgumentsMap(arguments map[string][]string) error {
	var err error
	argumentsMap := make(map[string]interface{})

	switch sm.Command {
	case "iface-stat":
		if value, ok := arguments["iface"]; ok {
			argumentsMap["iface"] = value[0]
			break
		} else {
			return &Error{Message: fmt.Sprintf("iface-stat should have an iface parameter: %+v", arguments)}
		}
	case "pcap-file":
		if value, ok := arguments["filename"]; ok {
			argumentsMap["filename"] = value[0]
		} else {
			return &Error{Message: fmt.Sprintf("pcap-file should have a filename parameter: %+v", arguments)}
		}

		if value, ok := arguments["output-dir"]; ok {
			argumentsMap["output-dir"] = value[0]
		} else {
			return &Error{Message: fmt.Sprintf("pcap-file should have a output-dir parameter: %+v", arguments)}
		}

		if value, ok := arguments["tenant"]; ok {
			argumentsMap["tenant"], err = strconv.Atoi(value[0])
			if err != nil {
				return &Error{Message: fmt.Sprintf("Tenant ID is not an intager: %+v", value)}
			}
		}

	case "conf-get":
		if value, ok := arguments["variable"]; ok {
			argumentsMap["variable"] = value[0]
		} else {
			return &Error{Message: fmt.Sprintf("conf-get should have a variable parameter: %+v", arguments)}
		}
	case "unregister-tenant-handler":
		return nil
	case "register-tenant-handler":
		return nil
	case "unregister-tenant":
		return nil
	case "register-tenant":
		return nil
	default:
		return &Error{Message: fmt.Sprintf("The command \"%s\" with %d arguments is unkown", sm.Command, len(arguments))}
	}
	sm.Arguments = &argumentsMap
	return nil
}

func (sm *SocketMessage) parseArgumentsList(arguments ...string) error {
	var err error
	argumentsMap := make(map[string]interface{})
	switch sm.Command {
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
		return &Error{Message: fmt.Sprintf("The command \"%s\" with %d arguments is unkown", sm.Command, len(arguments))}
	}
	sm.Arguments = &argumentsMap
	return nil
}
