package internal

import (
	"fmt"
	"io"

	pb "github.com/astronomical3/fewer_grpc/fewer"
)



//*****************************************************************************************
//  Definition of the FewerService struct that holds the Fewer Service, and whose methods
//    act as handlers for client method calls when an instance of this service is 
//    registered to a general gRPC server.
type FewerService struct {
	pb.UnimplementedFewerServiceServer
	serverLogger ServerLogger
}

// Constructor function for creating a new instance of the FewerService.
func NewFewerService(serverLogger ServerLogger) *FewerService {
	return &FewerService{serverLogger: serverLogger}
}

// Implementation of the GetAggregatesStream() RPC, which takes in NumberRequest messages,
//   adds 3 of them together, and every 3rd NumberRequest send, returns a sum (aggregate)
//   of the last 3 numbers in a NumberResponse.  Of course, this will be a very simple, 
//   incremental batch processing operation.
// This is an operation being used for testing whether or not it is possible to have the
//   server return back to clients FEWER responses than it receives requests (hence the 
//   name "Fewer Service").  If this operation is successful, it can be assumed that 
//   bidirectional streaming RPCs are useful for different types of batch processing.
func (s *FewerService) GetAggregatesStream(stream pb.FewerService_GetAggregatesStreamServer) error {
	i := 0
	sum := int32(0)
	for {
		// Try to receive a new NumberRequest, req, through the stream.
		req, err := stream.Recv()
		// If final request was already received from client...
		if err == io.EOF {
			if i % 3 != 0 {
				// Log final "leftover sum" into server log and return that sum to client if
				//   number of lefotver number requests is not 3.
				// Maybe this could be considered a "partial" operation, and could set off
				//   a warning.  We will simulate such a situation here...
				s.serverLogger.ServerLogWarn(
					"rpc",
					"pb.FewerService_GetAggregatesStream",
					fmt.Sprintf(
						"Leftover data not reported in last returned sum.  Actual final sum is %d.  Returning residual sum back to client...",
						sum,
					),
				)
				stream.Send(&pb.NumberResponse{Result: sum})
			} else {
				s.serverLogger.ServerLogInfo(
					"rpc",
					"pb.FewerService_GetAggregatesStream",
					"No leftover data after final sum.  Last sum returned is actual final sum.",
				)
			}
			return nil
		}

		// If receive error is some other non-nil error...
		if err != nil {
			s.serverLogger.ServerLogError(
				"rpc",
				"pb.FewerService_GetAggregatesStream",
				fmt.Sprintf("Could not receive latest request at iteration %d: %v", (i + 1), err),
			)
			return err
		}

		// If no receive error was received, or it is not the end of the stream of messages from the
		//   client...
		i++
		sum += req.InputNum
		s.serverLogger.ServerLogInfo(
			"rpc",
			"pb.FewerService_GetAggregatesStream",
			fmt.Sprintf("Received input number %d, sum is now %d", req.InputNum, sum),
		)
		if i % 3 == 0 {
			// Every 3 requests the service receives, it returns back the sum of those last 3 numbers
			//   received, and resets the sum back to 0.
			// If there is an error during the send, though, error is returned through gRPC runtime.
			s.serverLogger.ServerLogInfo(
				"rpc",
				"pb.FewerService_GetAggregatesStream",
				"3 input numbers have been added, sending back sum to client...",
			)
			if err := stream.Send(&pb.NumberResponse{Result: sum}); err != nil {
				s.serverLogger.ServerLogError(
					"rpc",
					"pb.FewerService_GetAggregatesStream",
					fmt.Sprintf("Could not send latest sum %d to client", sum),
				)
				return err
			}
			sum = int32(0)
		}
	}
}