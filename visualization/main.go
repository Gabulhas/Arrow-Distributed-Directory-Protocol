package main

import (
	"sync"
	"visualization/elements"
)

var Nodes sync.Map
var Links []elements.Connection
var AllUpdates []string

func init() {
	AllUpdates = make([]string, 0)
}

func main() {
	startServer()
}
