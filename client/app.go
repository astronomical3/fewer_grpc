package main

import (
	"log"

	"github.com/astronomical3/fewer_grpc/client/internal"
)

func main() {
	// Create a new CLI object that takes in the flags of this file's `go run`
	//   command, and then have it perform the Fewer Service's GetAggregatesStream()
	//   RPC to send over many number requests and receive back just a few 
	//   number responses.
	cliObj := internal.NewCli()
	cliObj.LoadAndParseFlags()
	if err := cliObj.PerformGetAggregatesOp(); err != nil {
		log.Printf("CLI object's PerformGetAggregatesOp operation ended in error: %v", err)
	}
}