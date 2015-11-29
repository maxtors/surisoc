package surisoc

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// Common variables for the suricta socket package
var (
	responseOK    = "OK"
	responseNotOK = "NOK"
	clientVersion = "0.1"
)

// SocketInit is the initial message sent to the suricata socket
type SocketInit struct {
	Version string `json:"version"`
}

// SuricataSocket is a socket connection object for interaction with suricata
type SuricataSocket struct {
	SocketPath    string
	Connection    net.Conn
	ValidCommands []string
}

// NewSuricataSocket creates a new suricata socket object
func NewSuricataSocket(path string) (*SuricataSocket, error) {
	suriSocket := &SuricataSocket{
		SocketPath: path,
	}

	// Try and connect to this socket
	err := suriSocket.Connect()
	if err != nil {
		return nil, err
	}

	// Perform the socket initialization routine
	err = suriSocket.InitConnect()
	if err != nil {
		return nil, err
	}

	// Get a list of valid commands
	response, err := suriSocket.Send("command-list")
	if err != nil {
		return nil, err
	}

	// If the response code was not "OK"
	if response.Return != responseOK {
		return nil, &Error{Message: fmt.Sprintf("Did not get OK response when finding valid commands: %+v", response)}
	}

	// Add the different commands to the valid commands list
	commands := response.Message.(map[string]interface{})["commands"]
	for _, command := range commands.([]interface{}) {
		suriSocket.ValidCommands = append(suriSocket.ValidCommands, command.(string))
	}

	// Return the newly created suricata socket object and nil as error
	return suriSocket, nil
}

// Connect creates the connection to the unix socket for a suricata socket
func (s *SuricataSocket) Connect() error {

	// Try and establish a connection to the socket
	conn, err := net.DialTimeout("unix", s.SocketPath, time.Second*3)
	if err != nil {
		return err
	}

	// Set the object connection variable and return nil
	s.Connection = conn
	return nil
}

// Close will close the current connection to the suricata unix socket
func (s *SuricataSocket) Close() error {
	return s.Connection.Close()
}

// InitConnect sets up a suricata socket connection
func (s *SuricataSocket) InitConnect() error {
	message := SocketInit{Version: clientVersion}
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Write the wanted data to the socket
	_, err = s.Connection.Write(bytes)
	if err != nil {
		return err
	}

	// Recive the response from the socket
	response, err := s.recive()
	if err != nil {
		return err
	}

	// If the response code in the JSON response is not "OK"
	if response.Return != responseOK {
		return &Error{fmt.Sprintf("Could not connect to Socket: %+v", response.Message)}
	}

	return nil
}

// Send will take a command and the given arguments and send them on the
// objects socket to the suricata process, it will check if the command
// is valid and can be performed
func (s SuricataSocket) Send(command string, arguments ...string) (*SocketResponse, error) {
	var err error

	// Check if this command is a valid command that we can send to the socket
	foundCommand := false

	// If the command is "command-list", then skipp checking of validity
	if command == "command-list" {
		foundCommand = true
	} else {
		for _, cmd := range s.ValidCommands {
			if command == cmd {
				foundCommand = true
				break
			}
		}
	}

	// If we did not find the command in the list of valid commands
	if !foundCommand {
		return nil, &Error{Message: fmt.Sprintf("Command %s is not valid", command)}
	}

	// The container for the message, either with or without arguments
	message := NewSocketMessage(command)
	err = message.ParseArgumentsList(arguments...)
	if err != nil {
		return nil, err
	}

	return s.SendMessage(message)
}

// SendMessage sends a Socket Message to the suricata socket
func (s SuricataSocket) SendMessage(msg *SocketMessage) (*SocketResponse, error) {
	// Marshal the socket message to a json []byte
	bytes, err := json.Marshal(msg)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Could not marshal socket message: %+v", msg)}
	}

	// Write the json []byte onto the socket
	_, err = s.Connection.Write(bytes)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Error while sending socket message: %s", err.Error())}
	}

	return s.recive()
}

// Recive will recive a response from the suricata socket
func (s *SuricataSocket) recive() (*SocketResponse, error) {

	time.Sleep(time.Millisecond * 10)

	// Read the data from the socket
	buffer := make([]byte, 8192)
	n, err := s.Connection.Read(buffer[:])
	if err != nil {
		return nil, err
	}

	// Parse the response and return to the caller
	response := SocketResponse{}
	json.Unmarshal(buffer[:n], &response)
	return &response, nil
}
