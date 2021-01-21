package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*

Apenas para debug

*/
//TODO: Remover

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
		selfNode.OutputState()
		break
	case cmd == "req": //cria request para o objeto
		selfNode.Request()
		break

	case cmd == "give_obj":
		selfNode.Obj = true
		break

	case cmd == "exit":
		os.Exit(1)
		break

	}

}
