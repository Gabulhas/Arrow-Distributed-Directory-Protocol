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
			fmt.Printf("%s", utils.StructToString(myChanReq))
		}

	}
}
