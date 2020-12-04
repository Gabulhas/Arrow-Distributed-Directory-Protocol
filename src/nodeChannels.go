package main

import (
	"fmt"
	"projeto/utils"
)

func ChanHandler() {
	for {
		select {
		case findReq := <-find:
			HandleFind(findReq)
			break
		case myChanReq := <-myChan:
			fmt.Printf("\nGot the object. %s", utils.StructToString(myChanReq))
			ReceiveObj(myChanReq)
		}

	}
}
