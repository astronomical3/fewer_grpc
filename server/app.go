package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/astronomical3/fewer_grpc/server/internal"
)

// Provide a name of the production-level or development-level server activity log file.
//   The filename is relative to the fewer_grpc/server/serverlogs/ directory path.
const serverLogProdFilename = "server.log"
const serverLogDevFilename = "server_devtest.log"

// Definition of the --address flag of the 'go run [fewer_grpc/server/]app.go' command.
var address = flag.String("address", "localhost", "address of server to serve on")

// Definition of the --port flag of the 'go run [fewer_grpc/server/]app.go' command.
var port = flag.Int("port", 50051, "port of server to serve on")

// Definition of the --prod flag of the 'go run [fewer_grpc/server/]app.go' command.
var prod = flag.Bool("prod", true, "indicates whether server is production server or development server")

func main() {
	// Load and parse the values of the flags provided in the 'go run' command.
	flag.Parse()
	
	// Create a listener, lis, that will listen to requests to the address and port
	//   provided by the 'go run' command flags.
	addrString := fmt.Sprintf("%s:%d", *address, *port)
	lis, err := net.Listen("tcp", addrString)
	if err != nil {
		log.Fatalf("fewer_grpc/server/app.go: net.Listen failed to listen: %v", err)
	}

	// Create a new GeneralFewerServer object.
	var serverLogFilename string
	if *prod {
		serverLogFilename = serverLogProdFilename
	} else {
		serverLogFilename = serverLogDevFilename
	}
	genServer := internal.NewGeneralFewerServer(serverLogFilename, lis, *prod)

	// Have the GeneralFewerServer object serve clients.  This method also handles
	//   shutdowns or server failures.
	genServer.ListenAndServe()
}