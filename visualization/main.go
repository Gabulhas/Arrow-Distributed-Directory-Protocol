package main

import "visualization/elements"

var Nodes map[string]elements.Node
var Links []elements.Link
var AllUpdates []elements.Node

func main() {

	initServer()
	startServer()
}

func initServer() {
	Nodes = make(map[string]elements.Node)
	AllUpdates = make([]elements.Node, 0)
}
