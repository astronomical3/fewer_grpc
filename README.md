# Example gRPC Fewer Service

Welcome to the gRPC Fewer Service!  This service is simply a demonstration on one way of how to organize gRPC client and server applications and the different components they would contain.  

## Wait, why did you call this the "Fewer Service"?

It's because the purpose of the service is to test whether gRPC bidirectional-streaming RPCs can be used for what I call "incremental batch processing".  That is, for every few inputs/requests that a service receives through a bidirectional stream, it sends an aggregate response for those few requests it gets.  In the case of this service, I decided to keep the core RPC operation, named `GetAggregatesStream()`, pretty simple -- for every 3 number requests it receives, it returns a sum of those 3 requests.  Note, too, that if the remaining last requests are not a set of 3 requests (that is, there are 1 or 2 requests left to process after the last 3 have been processed, because totalInputs is not a multiple of 3), a response containing a "residual sum" of the last 1 or 2 requests will be returned to the client right before the operation is finished.

## More about the application organization

I organized the client application to include a "core client object" that basically does the most essential actions for performing a gRPC operation: connect to the server, then perform a core gRPC operation.  Setting up and using the client application is as simple as:

1. Creating a client logging object (ClientLogger), depending on whether the client will be production- or development/test-grade:

   ```go
   clientLogFilename := "client.log"
   // If using production client object...
   clientLogger := NewClientLoggingObjectPROD(clientLogFilename)
   // If using client object...
   clientLogger := NewClientLoggingObjectDEV(clientLogFilename)
   ```
2. Calling the `NewCoreFewerSrvClient()` constructor function to create a new client application instance.  Be sure to use a `ClientLoggingObjectPROD` object if configuring the client to be production-grade (`isProd == true`).  Be sure to also defer the call of the client's `Close()` method to ensure server connections and client log files are properly closed, preventing resource leaks.

   ```go
   address := "some_hostname"
   port := 50051
   isProd := true
   coreClient := NewCoreFewerSrvClient(address, port, clientLogger, isProd)
   defer coreClient.Close()
   ```
3. Calling the client's `ConnectToServer()` method to make connection to the server application.

   ```go
   err := coreClient.ConnectToServer()
   if err != nil {
       // some code to handle the connection error, err
   }
   ```
4. Calling the client's `PerformGetAggregatesOp()` method to actually perform the core RPC operation.

   ```go
   totalInputs := 15
   err = coreClient.PerformGetAggregatesOp(totalInputs)
   if err != nil {
       // some code to handle the operation error, err
   }
   ```

Similarly, I organized the server application to include the "core service" (the actual Fewer Service and its `GetAggregatesStream()` RPC), as well as a "general server object" that gets a simple general gRPC server, so that this object can register its Fewer Service instance to that general gRPC server.  All that is needed to set up the server is: 

1. Calling the object's `NewGeneralFewerServer()` constructor function.
   ```go
   serverLogFilename := "server.log"
   lis := net.Listen("tcp", "some_hostname:50051")
   isProd := true
   genServer := NewGeneralFewerServer(serverLogFilename, lis, isProd)
   ```
2. Calling the server's `ListenAndServe()` method to start serving clients.  (NOTE: This method also serves to listen for any OS termination/interruption signals (ahem, **Ctrl+C**) so that it can gracefully shutdown the server.)
   ```go
   genServer.ListenAndServe()
   ```

Both the "core client object" and "general server" also have access to a "logging object", which logs a message simultaneously to 2 different log sources -- a terminal/stdout log, and a file log, without the need to rewrite the log message twice directly in the code for the "core client object".  I would just write it once using one of the logging object's methods, and from there let the logging object actually write out the full log message to the 2 different log sources.  I also created different logging objects based on whether I would simulate the client/server being in a production or development/test environment.

Additionally, I also decided to create an example client and server application, in CLI form, to demonstrate one way the "core client object" and "general server" application objects can be implemented.  To use the example applications I set up, simply clone this repo from GitHub.

## Using the example CLI applications

How to use the example client application (CLI): `
go run [fewer_grpc/client/]app.go [--address *hostname*] [--port *port_number*] [--prod={true|false}] [--totalInputs *num*]`

* `--address *hostname*`: Identify the hostname or address of the Fewer Service Server Application to connect to (default `"localhost"`).
* `--port *port_number*`: Identify the port of the Fewer Service Server Application to connect to (default `50051`).
* `--prod={true|false}`: Configure the Client Application to be either in a production environment (`true`) or development environment (`false`).  Default `true`.
* `--totalInputs *num*`: Specify the amount of numbers to send to the Fewer Service (default `15`).

How to use the example server application (CLI):
`go run [fewer_grpc/server/]app.go [--address *hostname*] [--port *port_number*] [--prod={true|false}]`

* `--address *hostname*`: Identify the address to serve the Fewer Service Server Application on (default `"localhost"`).
* `--port *port_number*`: Identify the port to serve the Fewer Service Server Application on (default `50051`).
* `--prod={true|false}`: Configure the Server Application to be in a production environment (true) or development environment (`false`).  Default `true`.

To shut down the Server App, you can just press **Ctrl+C**.

## Feedback

If you have comments, questions, etc., you can either:

* Raise an [issue](https://github.com/astronomical3/fewer_grpc/issues)
* Add a comment to the [Discussion board](https://github.com/astronomical3/fewer_grpc/discussions/3)
* Email me at astrobrunner@gmail.com
