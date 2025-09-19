package internal

import (
	"context"
	"fmt"
	"io"

	pb "github.com/astronomical3/fewer_grpc/fewer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

//***************************************************************************************************
// Definition of a core client object that can be easily set up and used in different implementations
//   of the Fewer Service Client Application (e.g., CLI, object included in a microservice).  This 
//   core client object is what actually dials up to the Fewer Service server and makes an RPC request
//   to the Fewer Service, and is received back some responses from the Fewer Service to log onto its
//   client activity logs.
type CoreFewerSrvClient struct {
	// Initial values upon creation of object
	addrString   string
	clientLogger ClientLogger
	isProd       bool

	// Obtained objects throughout connection and RPC execution process
	rpcCred      credentials.TransportCredentials
	grpcConn     *grpc.ClientConn
	grpcClient   pb.FewerServiceClient
}

// Constructor function for creating a new CoreFewerSrvClient that will dial up to the gRPC Fewer
//   Service server and perform operations from the service.
func NewCoreFewerSrvClient(address string, port int, clientLogger ClientLogger, isProd bool) *CoreFewerSrvClient {
	// Create the TCP address out of the given address/hostname and port.
	addrString := fmt.Sprintf("%s:%d", address, port)

	return &CoreFewerSrvClient{
		addrString:   addrString,
		clientLogger: clientLogger,
		isProd:       isProd,
	}
}

// Method of the CoreFewerSrvClient for dialing up to the gRPC Fewer Service server app and receiving
//   a client stub to the service.
func (c *CoreFewerSrvClient) ConnectToServer() error {
	// Get RPC credentials to use for dialing up to the Fewer Service server app.
	c.clientLogger.ClientLogInfo("method", "CoreFewerSrvClient.ConnectToServer", "Obtaining credentials for connecting core client to Fewer Service server...")
	if c.isProd {
		// NOTE: Right now, log a warning that this is not secure and that secure credentials for production
		//   environment clients need to be implemented.
		c.clientLogger.ClientLogWarn("method", "CoreFewerSrvClient.ConnectToServer", "Right now, insecure credentials will be used.  Working on implementing secure credentials for production environment operations...")
	}
	c.rpcCred = insecure.NewCredentials()

	// Dial up to the Fewer Service server app
	c.clientLogger.ClientLogInfo("method", "CoreFewerSrvClient.ConnectToServer", fmt.Sprintf("Connecting core client object to Fewer Service server at address %s...", c.addrString))

	var err error
	c.grpcConn, err = grpc.NewClient(c.addrString, grpc.WithTransportCredentials(c.rpcCred))
	if err != nil {
		c.clientLogger.ClientLogError("method", "CoreFewerSrvClient.ConnectToServer", fmt.Sprintf("Client failed to connect to server with grpc.NewClient: %v", err))
		return err
	}

	// Create a client stub to the Fewer Service.
	c.grpcClient = pb.NewFewerServiceClient(c.grpcConn)

	// An error of nil indicates that the client to the Fewer Service has been successfully made.
	return nil
}

// Method of the CoreFewerSrvClient that actually performs the operation of sending over to the Fewer Service server app
//   a bunch of pb.NumberRequest input messages, and receiving back pb.NumberResponse messages containing a sum of the
//   latest 3 inputs sent.
// This can be performed multiple times with the same client, by simply calling this function every time an operation is
//   requested.
func (c *CoreFewerSrvClient) PerformGetAggregatesOp(totalInputs int) error {
	// Create a done channel that will receive a close signal once the receiver
	//   goroutine in this operation has received all responses at end of operation.
	done := make(chan struct{})

	// Create a stream, numStream, through which the client will send NumberRequest
	//   messages to the Fewer Service Server through.
	numStream, err := c.grpcClient.GetAggregatesStream(context.Background())
	if err != nil {
		c.clientLogger.ClientLogError("method", "CoreFewerSrvClient.PerformGetAggregatesOp", fmt.Sprintf("Failure to open stream using GetAggregatesStream RPC: %v", err))
		return err
	}

	// Start up a sender goroutine that sends NumberRequest messages to the Fewer
	//   Service server via the opened numStream.
	go func() {
		for i := 1; i <= totalInputs; i++ {
			if err := numStream.Send(&pb.NumberRequest{InputNum: int32(i)}); err != nil {
				c.clientLogger.ClientLogWarn("method", "CoreFewerSrvClient.PerformGetAggregatesOp", fmt.Sprintf("Failed to send NumberRequest to server through numStream at request %d: %v", (i + 1), err))
				return
			}
		}
		numStream.CloseSend()
	}()

	// Start up a concurrent receiver goroutine that will receive some responses
	//   from the Fewer Service every 3 NumberRequest sends.
	// Also, create a recvErr channel that will return the error of the receive
	//   operation.
	recvErr := make(chan error)
	defer close(recvErr)
	go func() {
		for {
			resp, err := numStream.Recv()
			if err == io.EOF {
				// If last response was already received...
				c.clientLogger.ClientLogInfo("method", "CoreFewerSrvClient.PerformGetAggregatesOp", "Received all responses, closing done channel...")
				close(done)
				recvErr <- nil
				return
			}
			if err != nil {
				// If error results during a receive...
				c.clientLogger.ClientLogError("method", "CoreFewerSrvClient.PerformGetAggreatesOp", fmt.Sprintf("Failed to receive a response: %v", err))
				close(done)
				recvErr <- err
				return
			}
			c.clientLogger.ClientLogInfo("method", "CoreFewerSrvClient.PerformGetAggregatesOp", fmt.Sprintf("Received response from Fewer Service server: %v", resp))
		}
	}()

	// Block until all sending and receiving is finished.
	<-done
	
	// Get the final error of the receive operation from the receiving goroutine.
	//   If the final error was not nil, return the actual error. Otherwise, the
	//   method returns nil, indicating the entire operation was successful.
	finalErr := <-recvErr
	if finalErr != nil {
		return finalErr
	}

	return nil
}

// Method of the CoreFewerSrvClient for closing the client's resources (the gRPC 
//   connection, the client log file used by the attached clientLogger, etc.)
func (c *CoreFewerSrvClient) Close() {
	c.grpcConn.Close()
	c.clientLogger.Close()
}