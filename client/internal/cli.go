package internal

import (
	"flag"
)



//*************************************************************************************************
// Definition of a CLI object that takes in the flags to the `go run` command that runs the client
//   main setup code.
type Cli struct {
	address     *string
	port        *int
	totalInputs *int
	prod        *bool
}

// Constructor function that creates a new instance of the *Cli object
func NewCli() *Cli {
	return &Cli{}
}

// Method of the Cli object that loads and parses the flags to the client main setup's
//   `go run` command.
func (cli *Cli) LoadAndParseFlags() {
	// Address and port of the Fewer Service server to connect to
	cli.address = flag.String("address", "localhost", "address or hostname of the Fewer Service server to connect to")
	cli.port = flag.Int("port", 50051, "port of the Fewer Service server to connect to")

	// Maximum number of requests to send to the service
	cli.totalInputs = flag.Int("totalInputs", 15, "maximum number of requests to send to Fewer Service")
	
	// Whether the client is production-grade or not
	cli.prod = flag.Bool("prod", true, "indicates whether the client is a production (true) or development/test (false) client")

	flag.Parse()
}

// Method of the Cli object that creates a core client object and performs the process 
//   of sending over number request messages to the Fewer Service server, in an attempt
//   to get a few number responses back.
func (cli *Cli) PerformGetAggregatesOp() error {
	// Create a ClientLogger, based on whether the client will be production
	//   or development/testing (non-production).
	const clientLogProdFilename = "client.log"
	const clientLogDevFilename = "client_devtest.log"
	var clientLogger ClientLogger
	if *cli.prod {
		clientLogger = NewClientLoggingObjectPROD(clientLogProdFilename)
	} else {
		clientLogger = NewClientLoggingObjectDEV(clientLogDevFilename)
	}

	// Create core client object.
	coreClient := NewCoreFewerSrvClient(*cli.address, *cli.port, clientLogger, *cli.prod)
	defer coreClient.Close()

	// Connect the core client to the Fewer Service server.
	err := coreClient.ConnectToServer()
	if err != nil {
		return err
	}

	// Perform the Fewer Service's bidirectional-streaming GetAggregatesStream() RPC,
	//   sending over many number requests, and receiving back only a few number responses
	//   from the Fewer Service.
	err = coreClient.PerformGetAggregatesOp(*cli.totalInputs)
	if err != nil {
		return err
	}

	// A nil error indicates successful connection and operation
	return nil
}