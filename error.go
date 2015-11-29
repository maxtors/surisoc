package surisoc

// Error is a custom error struct for the SuriSock package
type Error struct {
	Message string
}

// Error gives back the error message
func (e *Error) Error() string {
	return e.Message
}
