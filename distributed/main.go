package main

import (
	"os"
	"tour/distributed/nodes"
)

var workerPorts = []string{"4000"}

func main() {

	switch os.Args[1] {
	case "master":
		nodes.StartServer()
	case "worker":
		nodes.StartWorkers(workerPorts)
		//case "webserver":

	}

}
