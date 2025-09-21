package internal

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/astronomical3/fewer_grpc/fewer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//*************************************************************************
// Definition of the general gRPC server that will host the Fewer Service.
type GeneralFewerServer struct {
	// Added upon creation of the server
	listener     net.Listener

	// Received later on in the setup.
	grpcServer   *grpc.Server
	serverLogger ServerLogger
	srv          *FewerService
}

// Create a new general gRPC server, and create a new server logging object depending on whether the server 
//   will be production or development/test.
func NewGeneralFewerServer(serverLogFilename string, lis net.Listener, isProd bool) *GeneralFewerServer {
	// Obtain a new general gRPC server
	grpcServer := grpc.NewServer()

	// Obtain a new server logging object depending on whether the server will be production- or 
	//   development/test-grade.
	var serverLogger ServerLogger
	if isProd {
		serverLogger = NewServerLoggingObjectPROD(serverLogFilename)
	} else {
		serverLogger = NewServerLoggingObjectDEV(serverLogFilename)
	}

	// Create a new instance of the Fewer Service.
	srv := NewFewerService(serverLogger)

	return &GeneralFewerServer{
		listener:     lis,
		grpcServer:   grpcServer,
		serverLogger: serverLogger,
		srv:          srv,
	}
}

// Method of the GeneralFewerServer that is used for registering its new Fewer Service instance to its general
//   gRPC server, registering a gRPC reflection service (for providing RPC and other useful information about
//   the Fewer Service to tools such as grpcurl), and then setting up a channel to listen to OS termination or
//   interruption signals and a goroutine for serving Fewer Service clients.
func (fs *GeneralFewerServer) ListenAndServe() {
	// Register the Fewer Service instance and an instance of the gRPC reflection service to the gRPC server.
	pb.RegisterFewerServiceServer(fs.grpcServer, fs.srv)
	reflection.Register(fs.grpcServer)

	// Create a channel, sigChan, that will listen to an OS termination or interruption signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Have the server listen to all FewerService-specific requests.
	go func() {
		fs.serverLogger.ServerLogInfo(
			"method",
			"GeneralFewerServer_ListenAndServe",
			fmt.Sprintf("FewerServer listening on address %v", fs.listener.Addr()),
		)
		if err := fs.grpcServer.Serve(fs.listener); err != nil {
			fs.serverLogger.ServerLogError(
				"method",
				"GeneralFewerServer_ListenAndServe",
				fmt.Sprintf("Failed to serve via grpcServer.Serve() method: %v", err),
			)
			fs.serverLogger.ServerLogError(
				"method",
				"GeneralFewerServer_ListenAndServe",
				"Closing server log, returning exit code 1...",
			)
			fs.serverLogger.Close()
			fs.listener.Close()
			os.Exit(1)
		}
	}()

	// Block until an interruption/termination signal is received.
	sig := <-sigChan
	fs.serverLogger.ServerLogInfo(
		"method",
		"GeneralFewerServer_ListenAndServe",
		fmt.Sprintf("Received signal (%s), starting graceful shutdown...", sig.String()),
	)
	fs.shutdown()
}

// Internal method of the GeneralFewerServer for ensuring graceful stop of
//  gRPC server when an OS termination/interruption signal is issued.
func (fs *GeneralFewerServer) shutdown() {
	fs.grpcServer.GracefulStop()
	fs.serverLogger.ServerLogInfo(
		"method",
		"GeneralFewerServer_Shutdown",
		"gRPC server gracefully stopped.",
	)
	fs.serverLogger.Close()
}