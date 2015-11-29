# SuriSoc
A Go-lang package for interaction with the suricata command unix socket

## Installation
- Have golang installed
- >> go get github.com/Maxtors/suricon

## Usage
```go
// Create a new Suricata Socket session
session, err = surisoc.NewSuricataSocket(socketPath)
if err != nil {
    log.Fatalf("Error: %s\n", err.Error())
}
defer session.Close()

// Send the "version" command to get the version of Suricata that is running
response, err = session.Send("version")
if err != nil {
    log.Fatalf("Error: %s\n", err.Error())
}

// Convert the response.Message to a string
res, err := response.ToString()
if err != nil {
    log.Fatalf("Error: %s\n", err.Error())
}

// Print results
fmt.Println(res)
```
