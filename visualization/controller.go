package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"visualization/elements"
)

//TODO: Mudar para package
//TODO: dividir as funções em vários ficheiros

var re = regexp.MustCompile(`http://|/find|/myChan`)
var requestHistory []string
var queueHistory []string
var ownerHistory []string
var Queues [][]elements.Node
var QueuesMutex = &sync.RWMutex{}
var currentOwner = ""
var totalQueuesPrinted = 0

//TODO: mudar de Mutex para RW mutex

func startServer() {

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	r.HandleFunc("/", root).Methods("GET")
	r.HandleFunc("/data", data).Methods("GET")
	r.HandleFunc("/queue", queue).Methods("GET")
	r.HandleFunc("/updateState", updateState).Methods("POST")
	// TODO: remover /logs, visto que não apresentam a ordem correta de chegada
	r.HandleFunc("/logs", getLogs).Methods("GET")
	r.HandleFunc("/requestAll", requestAll).Methods("GET")

	log.Fatal(http.ListenAndServe(os.Getenv("address"), r))
}

func root(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("./assets/html/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusFound)
	}
	if err := tmpl.ExecuteTemplate(w, tmpl.Name(), nil); err != nil {
		log.Fatalf("homeHandler: %+v", err)
	}
}

func data(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	response := new(elements.VisResponse)

	var tempNodes []elements.Node
	var tempLinks []elements.Connection
	var tempQueueCons []elements.Connection

	Nodes.Range(func(key, value interface{}) bool {
		v := value.(elements.Node)
		tempNodes = append(tempNodes, v)

		if v.Link != "" {
			tempLinks = append(tempLinks, elements.Connection{
				Source: v.MyAddress,
				Target: v.Link,
			})
		}
		if v.WaiterChan != "" {
			tempQueueCons = append(tempQueueCons, elements.Connection{
				Source: v.MyAddress,
				Target: v.WaiterChan,
			})
		}

		return true
	},
	)

	response.Nodes = tempNodes
	response.Links = tempLinks
	response.QueueCons = tempQueueCons

	json.NewEncoder(w).Encode(response)
}

func updateState(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var update elements.Node
	_ = json.NewDecoder(r.Body).Decode(&update)

	update.Link = re.ReplaceAllString(update.Link, ``)
	update.WaiterChan = re.ReplaceAllString(update.WaiterChan, ``)

	valueInterface, _ := Nodes.Load(update.MyAddress)

	Nodes.Store(update.MyAddress, update)

	json.NewEncoder(w).Encode("Successful")

	if update.Type == 2 {
		return
	}
	QueuesMutex.Lock()
	updateQueues(update, valueInterface)
	QueuesMutex.Unlock()

	QueuesMutex.RLock()
	totalQueuesPrinted = totalQueuesPrinted + 1
	PrettyPrintQueue(Queues)
	AllUpdates = append(AllUpdates, OtherPrintQueue(Queues))
	QueuesMutex.RUnlock()
}

func OtherPrintQueue(queues [][]elements.Node) string {
	finalString := fmt.Sprintf("%d-", totalQueuesPrinted)
	for _, tempQueue := range queues {
		for _, node := range tempQueue {
			finalString = finalString + fmt.Sprintf("%s ", node.MyAddress)
		}
		finalString = finalString + " # "
	}
	return finalString

}

func PrettyPrintQueue(queues [][]elements.Node) {
	finalString := fmt.Sprintf("%d", totalQueuesPrinted)
	for _, tempQueue := range queues {

		switch len(tempQueue) {
		case 0:
			finalString = finalString + fmt.Sprintf("ERROR")
			break
		case 1:
			finalString = finalString + fmt.Sprintf("[%s]", tempQueue[0].MyAddress)
			break
		default:
			finalString = finalString + fmt.Sprintf("[%s -> %s]", tempQueue[0].MyAddress, tempQueue[len(tempQueue)-1].MyAddress)
			break
		}
	}
	fmt.Println(finalString)
}


func nodeChange() {
	for {
		select {

		}
	}
}

func updateQueues(update elements.Node, valueInterface interface{}) {

	previous, ok := valueInterface.(elements.Node)

	switch update.Type {
	case 0:
		if !ok {
			return
		}

		if previous.Type == 3 {
			for i, queue := range Queues {

				if queue[0].MyAddress == update.MyAddress {
					Queues[i] = Queues[i][1:]
					if len(Queues[i]) == 0 {
						Queues = removeFromQueue(Queues, i)

					}
				}
			}

		}
		if previous.Type == 1 {
			for i, queue := range Queues {
				if queue[0].MyAddress == update.MyAddress {
					temp := queue[1:]
					Queues[i] = Queues[0]
					Queues[0] = temp
					return
				}
			}
		}
		break
	case 1:
		for i, queue := range Queues {
			if queue[0].MyAddress == update.MyAddress {
				if len(queue) == 1 {
					Queues = removeFromQueue(Queues, i)
				} else {
					Queues[i] = Queues[i][1:]
				}
			}
		}
		break
	case 3:
		if !ok {
			return
		}
		if previous.Type == 3 {
			return
		}

		firstQueue := -1
		secondQueue := -1

		for i, queue := range Queues {
			if queue[len(queue)-1].MyAddress == update.MyAddress {
				firstQueue = i
			}
			if queue[0].MyAddress == update.WaiterChan {
				secondQueue = i
			}
		}

		if firstQueue == -1 || secondQueue == -1 {

			var temp elements.Node
			foundSecond := false

			if valueInterface, ok := Nodes.Load(update.MyAddress); ok {
				temp, _ = valueInterface.(elements.Node)
				foundSecond = true
			}
			if firstQueue == -1 {
				for i, queue := range Queues {
					if queue[len(queue)-1].WaiterChan == update.MyAddress {
						Queues[i] = append(Queues[i], update)
						if foundSecond {
							Queues[i] = append(Queues[i], temp)
							foundSecond = false
						}
					}
				}
			}
			if secondQueue == -1 && foundSecond && firstQueue != -1{
				Queues[firstQueue] = append(Queues[firstQueue], temp)
			}

			fmt.Printf("IMPOSSÍVEL %d:%s %d:%s\n", firstQueue, update.WaiterChan, secondQueue, update.MyAddress)
			return
		}
		Queues[firstQueue] = append(Queues[firstQueue], Queues[secondQueue]...)
		Queues = removeFromQueue(Queues, secondQueue)

		break
	case 4:
		for i, queue := range Queues {
			if queue[len(queue)-1].WaiterChan == update.MyAddress {
				Queues[i] = append(Queues[i], update)
				return
			}
		}
		Queues = append(Queues, []elements.Node{update})
		break
	}

}

func removeFromQueue(queue [][]elements.Node, i int) [][]elements.Node {
	return append(queue[:i], queue[i+1:]...)
}

func queue(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	response := new(elements.QueueResponse)

	//var currentNode elements.Node
	//var nextNode elements.Node
	var QueueNodesAddresses []string

	QueuesMutex.RLock()
	if len(Queues) > 0 {
		for _, node := range Queues[0] {
			QueueNodesAddresses = append(QueueNodesAddresses, node.MyAddress)
		}
	}
	QueuesMutex.RUnlock()
	response.QueueNodes = QueueNodesAddresses
	response.Requesting = requestHistory
	/*
		requestHistory = nil

		Nodes.Range(func(key, value interface{}) bool {
			node := value.(elements.Node)
			if node.Type < 2 {
				currentOwner = node.MyAddress
			}
			return true
		},
		)

	*/
	/*
		for _, node := range Nodes {
			if node.Type < 2 {
				currentOwner = node.MyAddress
			}
		}
	*/

	/*
		response.OwnerHistory = ownerHistory
		response.CurrentOwner = currentOwner
		tempStruct, _ := Nodes.Load(currentOwner)

		currentNode = tempStruct.(elements.Node)

		//fmt.Println("currentNode", currentNode)
		if currentOwner == "" {
			return
		}

		for currentNode.WaiterChan != "" {
			tempStruct, _ := Nodes.Load(re.ReplaceAllString(currentNode.WaiterChan, ``))
			nextNode = tempStruct.(elements.Node)

			response.QueueNodes = append(response.QueueNodes, nextNode.MyAddress)
			currentNode = nextNode
		}

		//dever haver algoritmo mais simples que este
		pivot := 0
		flag := true
		for i := 0; i < len(queueHistory); i++ {
			if pivot < len(response.QueueNodes) {
				if queueHistory[i] == response.QueueNodes[pivot] {
					pivot = i
					flag = false
					break
				}
			}
		}

		startPoint := len(queueHistory) - pivot

		if flag {
			startPoint = 0
		}

		for i := startPoint; i < len(response.QueueNodes); i++ {
			response.QueueHistory = append(response.QueueHistory, response.QueueNodes[i])
		}
		/*
		fmt.Println(queueHistory)
		fmt.Println(response.QueueNodes)
		fmt.Println("------------------------------------")
				ownerHistory = nil
				queueHistory = response.QueueNodes
	*/

	json.NewEncoder(w).Encode(response)
}

func getLogs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "text/plain")
	var sb strings.Builder
	sb.Grow(len(AllUpdates))
	for _, logStamp := range AllUpdates {
		fmt.Fprintf(&sb, "%s\n", logStamp)
	}
	w.Write([]byte(sb.String()))

}

func requestAll(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var wg sync.WaitGroup

	w.Header().Set("Content-Type", "text/plain")

	Nodes.Range(func(key, value interface{}) bool {
		element := value.(elements.Node)
		wg.Add(1)
		go remoteRequest(element.MyAddress, &wg)
		return true
	},
	)
	wg.Wait()
	w.Write([]byte("Successful"))
}

// TODO: usar esta função para o remote request de um só node, do "double click"
func remoteRequest(address string, wg *sync.WaitGroup) {

	defer (*wg).Done()

	_, err := http.Get(fmt.Sprintf("http://%s/remoteRequest", address))
	if err != nil {
		fmt.Println(err)
	}
}
