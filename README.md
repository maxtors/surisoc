# SuriCon
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

response, err = session.Send("version")
if err != nil {
    log.Fatalf("Error: %s\n", err.Error())
}

res, err := response.ToString()
if err != nil {
    log.Fatalf("Error: %s\n", err.Error())
}
fmt.Println(res)
```
