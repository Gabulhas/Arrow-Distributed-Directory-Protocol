package main

import (
	"bufio"
	"fmt"
	"os"
	"projeto/Channels"
	"strings"
)

func ShellStart() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nnode> ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		runCommand(cmdString)
	}
}

func runCommand(command string) {
	args := strings.Split(command, " ")
	cmd := strings.TrimSpace(args[0])

	switch {
	case cmd == "state":
		OutputState()
		break
	case cmd == "req": //cria request para o objeto
		makeRequest()
		break

	case cmd == "give_obj":
		Node.Obj = true
		break

	case cmd == "exit":
		os.Exit(1)
		break

	}

}

//Mudar esta função de ficheiro
func makeRequest() {
	accessRequest := Channels.AccessRequest{
		GiveAccess: Channels.GiveAccess{
			WaiterChan: Node.MyChan,
		},
		Link: Node.Find,
	}
	go SendThroughLink(accessRequest)

}
