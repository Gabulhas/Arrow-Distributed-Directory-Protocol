
# Arrow Distributed Directory Protocol


This project involves the study and implementation of a protocol for distributed systems, the [Arrow Distributed Directory Protocol](http://cs.brown.edu/people/mph/DemmerH98/disc.pdf).

The realization of this comprises the topics of programming in Go Language, Data Structures, concurrency, distributed systems and algorithms and their visualization.

The elaboration of this project took place in the end of project unit
Degree course in Computer Engineering.

It also includes an essay where I explain the Protocol and the process of implementation I went through.


Check [#Explanation](/#Explanation)

## Coordinators 

- [Professor Simão Melo de Sousa](https://scholar.google.pt/citations?user=nuuKV9cAAAAJ&hl=pt-PT)
- [Professor Hugo Torres Vieira](https://scholar.google.com/citations?user=Y5yb7XEAAAAJ&hl=en)
  
## Algorithm Explanation
(TODO ;) )

## Features

- **Implementation** of the Protocol in Go
    - Nodes communicate using HTTP.

- **Live Visualization** of the System in a Web App
    - Implementation of a Visualization Node that receives information from each Node.

- **Event Logging** to prove the correction of the System
    - Request History - order which the Visualization Node process each update of a request.
    - Owner History - order which the Nodes became Owners of the Mobile Object
    - Queue History - order which the Nodes got into the "main" queue/queue that directly follows the current Owner.


- Remotely control the System Nodes through the Visualization
- **Deployment** of a whole system (and visualization) with **Docker**
  
## Demos (on YouTube)

- [System Deployment (using Docker)](https://youtu.be/D1UuCL_JIBI)
- [Interface Showcase](https://youtu.be/de_02kPA3qc)
- [Event Logging](https://youtu.be/FIdTcE4UJsg)

  
## Run Locally

Cloning/Downloading this project
```bash
git clone https://github.com/Gabulhas/projeto-relatorio
```


### Builing Normally

Building the Node (Change Directory to "src")
```bash
go build -o Node .
```

Executing one Node 
```bash
#Arguments:
#  -address string (Required)
#    	Node's Address 
#
#  -link string (if any)
#    	address of the Node to which it's connected
#
#  -requests 
#    	If this Node, when Idle, preforms Object Requests (default true)
#
#  -type int (Required)
#    		Owner with Request  - 0
#	        Owner Terminal      - 1
#	        Idle	            - 2
#	        Waiter with Request - 3
#	        Waiter Terminal     - 4
#
#  -visualization string (if any)
#    	Address of the Visualization Node 


# Example, a Node that
# is an Idle Node (2)
# which address is 127.0.0.2
# is connected to 127.0.0.3
# where the visualization Node has the address 127.0.0.10

./Node -type 2 -address 127.0.0.2 -link 127.0.0.2 -visualization 127.0.0.10
```





### Using Docker/Makefile
Clone the project


  
## Screenshots

#### Interface Overview
- Each Circle in the graph is a different machine, and the arrows between them indicate their connections.
- The current "main" queue and the object owners are displayed on the right.
- The events are logged and displayed at the bottom tables.

![Interface Overview](https://github.com/Gabulhas/projeto-relatorio/raw/master/relatorio_overview.png)

  
## Lessons Learned

These were the topics which I had the opportunity to learn more of/new experiences:
- Distributed Systems
- Concurrency
- Go Programming Language
- Academic Writting/Formal Writting
- Working with Professors/Academicians (which I'm grateful to do so)

## Acknowledgements/Related Projects

 - [The Arrow Distributed Directory Protocol](http://cs.brown.edu/people/mph/DemmerH98/disc.pdf) - The Original Paper of the Algorithm/Protocol which I followed for the implementation.
 - [Bully Algorithm by @TimTosi](https://github.com/TimTosi/bully-algorithm) - Where I got the inspiration for the Visualization 
 - [Arvy by @Infinisil](https://github.com/Infinisil/arvy)

  
  

