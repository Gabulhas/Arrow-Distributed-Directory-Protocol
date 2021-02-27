package main

import (
	"sync"
	"visualization/elements"
)

var Nodes sync.Map
var Links []elements.Connection
var AllUpdates []elements.Node

func init(){
	AllUpdates = make([]elements.Node, 0)
}

func main() {
	startServer()
}