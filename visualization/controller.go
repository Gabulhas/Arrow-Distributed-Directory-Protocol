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
var currentOwner elements.Node
var QueuesMutex = &sync.RWMutex{}
var ChangeChannel chan elements.Node

//Lookup table
var NodesInQueue map[string]int

//Debug
var totalQueuesPrinted = 0

func init() {
	ChangeChannel = make(chan elements.Node, 40)
	NodesInQueue = make(map[string]int)
	go nodeChange()
}

func startServer() {

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	r.HandleFunc("/", root).Methods("GET")
	r.HandleFunc("/data", data).Methods("GET")
	r.HandleFunc("/queue", queue).Methods("GET")
	r.HandleFunc("/updateState", updateState).Methods("POST")

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

	Nodes.Store(update.MyAddress, update)

	json.NewEncoder(w).Encode("Successful")

	if update.Type != 2 {
		ChangeChannel <- update
	}
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
		case newChange := <-ChangeChannel:
			QueuesMutex.Lock()
			updateQueues(newChange)
			QueuesMutex.Unlock()

			//TODO: apenas para debug
			QueuesMutex.RLock()
			totalQueuesPrinted = totalQueuesPrinted + 1
			PrettyPrintQueue(Queues)
			AllUpdates = append(AllUpdates, OtherPrintQueue(Queues))
			QueuesMutex.RUnlock()
			break
		}
	}
}

//TODO: Limpar o código e dividir entre controller, queue e graph (ficheiros)

func updateQueues(update elements.Node) {

	fmt.Printf("(%d, %s -> %s)", update.Type, update.MyAddress, update.WaiterChan)

	switch update.Type {
	case 4:
		if _, isIn := NodesInQueue[update.MyAddress]; isIn {
			return
		}
		NodesInQueue[update.MyAddress] = update.Type
		for i, queue := range Queues {
			if len(queue) == 0 {
				continue
			}
			if queue[len(queue)-1].WaiterChan == update.MyAddress {
				Queues[i] = append(Queues[i], update)
				return
			}
		}
		Queues = append(Queues, []elements.Node{update})

		break
	case 3:
		if tipo, ok := NodesInQueue[update.MyAddress]; ok && tipo == 3 {
			return
		}

		_, IsNodeAIn := NodesInQueue[update.MyAddress]
		typeValue, IsNodeBIn := NodesInQueue[update.WaiterChan]

		NodesInQueue[update.MyAddress] = update.Type
		NodesInQueue[update.WaiterChan] = typeValue
		fmt.Printf("A: %t, B:%t", IsNodeAIn, IsNodeBIn)

		NodesInQueue[update.MyAddress] = update.Type

		nextInterface, _ := Nodes.Load(update.WaiterChan)
		nextNode := nextInterface.(elements.Node)

		//TODO: Mudar IFS

		if IsNodeAIn && IsNodeBIn {
			firstQueue := -1
			secondQueue := -1
			for i, queue := range Queues {
				if queue[len(queue)-1].MyAddress == update.MyAddress {
					firstQueue = i
				} else if queue[0].MyAddress == update.WaiterChan {
					secondQueue = i
				}
			}
			if secondQueue == -1 || firstQueue == -1 {
				return
			}
			Queues[firstQueue] = append(Queues[firstQueue], Queues[secondQueue]...)
			removeFromQueues(secondQueue)
			return

		}

		if IsNodeAIn && !IsNodeBIn {
			for i, queue := range Queues {
				if queue[len(queue)-1].MyAddress == update.MyAddress {
					Queues[i] = append(Queues[i], nextNode)
					return
				}
			}
		}

		if !IsNodeAIn && IsNodeBIn {
			for i, queue := range Queues {
				if queue[0].MyAddress == update.WaiterChan {
					Queues[i] = append([]elements.Node{update}, Queues[i]...)
					return
				}
			}
		}

		if !IsNodeAIn && !IsNodeBIn {

			currentQueue := []elements.Node{update, nextNode}

			currentQueueLocation := -1

			for i, queue := range Queues {

				//TODO: limpar estes IFs
				if queue[len(queue)-1].WaiterChan == update.MyAddress {
					if currentQueueLocation == -1 {
						Queues[i] = append(Queues[i], currentQueue...)
						currentQueueLocation = i
					} else {
						Queues[i] = append(Queues[i], Queues[currentQueueLocation]...)
						removeFromQueues(currentQueueLocation)
					}
				} else if nextNode.WaiterChan == queue[0].MyAddress {
					if currentQueueLocation == -1 {
						Queues[i] = append(currentQueue, Queues[i]...)
						currentQueueLocation = i
					} else {
						Queues[currentQueueLocation] = append(Queues[currentQueueLocation], Queues[i]...)
						removeFromQueues(i)
					}
				}

			}
			if currentQueueLocation == -1 {
				Queues = append(Queues, currentQueue)
			}
		}

		break

	case 1:
		foundNodeQueue := -1
		currentOwner = update
		for i, queue := range Queues {
			if queue[0].MyAddress == update.MyAddress {
				foundNodeQueue = i
			}
		}
		if foundNodeQueue == -1 {
			return
		}

		if len(Queues[foundNodeQueue]) == 1 {
			removeFromQueues(foundNodeQueue)
		} else {
			Queues[foundNodeQueue] = Queues[foundNodeQueue][1:]
		}
		delete(NodesInQueue, update.MyAddress)

		break
	case 0:

		delete(NodesInQueue, update.MyAddress)
		currentOwner = update
		for i, queue := range Queues {
			if queue[0].MyAddress == update.MyAddress {
				if len(queue) > 1 {
					Queues[i] = Queues[i][1:]
					temp := queue[1:]
					Queues[i] = Queues[0]
					Queues[0] = temp
					return
				} else {
					removeFromQueues(i)
				}
			} else if queue[0].MyAddress == update.WaiterChan {
				temp := queue
				Queues[i] = Queues[0]
				Queues[0] = temp
			}
		}

		fmt.Println("Impossible case, OWR not (pointing to) head")
		break
	}

}

func removeFromQueues(i int) {
	Queues = append(Queues[:i], Queues[i+1:]...)
}

func queue(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	response := new(elements.QueueResponse)

	//var currentNode elements.Node
	//var nextNode elements.Node
	var QueueNodesAddresses []string

	QueuesMutex.RLock()

	if len(Queues) > 0 && Queues[0][0].MyAddress == currentOwner.WaiterChan {
		for _, node := range Queues[0] {
			QueueNodesAddresses = append(QueueNodesAddresses, node.MyAddress)
		}
	}
	QueuesMutex.RUnlock()
	response.QueueNodes = QueueNodesAddresses
	response.Requesting = requestHistory
	response.CurrentOwner = currentOwner.MyAddress
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
