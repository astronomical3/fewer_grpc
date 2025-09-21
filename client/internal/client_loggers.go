package internal

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

//********************************************************************************************
// Definition of a ClientLogger interface that includes common logging operation methods that
//   will be used across multiple concrete logger types.
type ClientLogger interface {
	ClientLogInfo(key, value, message string)
	ClientLogWarn(key, value, message string)
	ClientLogError(key, value, message string)
	Close()
}



//*********************************************************************************************
// Definition of a production-level client activity logger that logs both on the 
//   terminal/standard output and a file.
type ClientLoggingObjectPROD struct {
	// Log only current client session activity on standard output
	terminalLogger   log.Logger
	// Log all client sessions' history on a file
	clientFileLogger log.Logger
	// Log file that the clientFileLogger operates on.  Stored in this field 
	//   so that it can be properly closed when client app is exited.
	clientLogFile    *os.File
}

// Constructor function that creates a production-level logger for logging all client activity,
//   INFO-level and above.
func NewClientLoggingObjectPROD(clientLogFilename string) *ClientLoggingObjectPROD {
	// Create terminal logger
	terminalLogger := log.NewLogfmtLogger(os.Stdout)
	terminalLogger = level.NewFilter(terminalLogger, level.AllowInfo())
	terminalLogger = log.With(terminalLogger, "time", log.DefaultTimestampUTC)

	// Add filepath to clientlogs directory to get full client log file path.
	//   clientLogFilename = "github.com/astronomical3/fewer_grpc/client/clientlogs/production/" + clientLogFilename

	// Create or open the client log file object
	clientLogFile, err := os.OpenFile(clientLogFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	// Create the logger to operate logging operations on the client log file
	clientFileLogger := log.NewLogfmtLogger(clientLogFile)
	clientFileLogger = level.NewFilter(clientFileLogger, level.AllowInfo())
	clientFileLogger = log.With(clientFileLogger, "time", log.DefaultTimestampUTC)

	return &ClientLoggingObjectPROD{
		terminalLogger:   terminalLogger,
		clientFileLogger: clientFileLogger,
		clientLogFile:    clientLogFile,
	}
}

// Method of the ClientLoggingObjectPROD that is used for simultaneously logging INFO-level
//   activity on both the terminal and client log file.
func (clop *ClientLoggingObjectPROD) ClientLogInfo(key, value, message string) {
	go level.Info(clop.terminalLogger).Log(key, value, "message", message)
	level.Info(clop.clientFileLogger).Log(key, value, "message", message)
}

// Method of the ClientLoggingObjectPROD that is used for simultaneously logging WARN-level
//   activity on both the terminal and client log file.
func (clop *ClientLoggingObjectPROD) ClientLogWarn(key, value, message string) {
	go level.Warn(clop.terminalLogger).Log(key, value, "message", message)
	level.Warn(clop.clientFileLogger).Log(key, value, "message", message)
}

// Method of the ClientLoggingObjectPROD that is used for simultaneously logging ERROR-level
//   activity on both the terminal and client log file.
func (clop *ClientLoggingObjectPROD) ClientLogError(key, value, message string) {
	go level.Error(clop.terminalLogger).Log(key, value, "error", message)
	level.Error(clop.clientFileLogger).Log(key, value, "error", message)
}

// Method of the ClientLoggingObjectPROD that is used for closing the client log file properly.
func (clop *ClientLoggingObjectPROD) Close() {
	clop.clientLogFile.Close()
}



//***********************************************************************************************************
// Definition of a development-level client activity logger that logs both on the 
//   terminal/standard output and a file.
type ClientLoggingObjectDEV struct {
	// Inherit the ClientLoggingObjectPROD methods, because they satify the ClientLogger
	//   interface
	ClientLoggingObjectPROD
	// Log only current session activity on standard output
	terminalLogger   log.Logger
	// Log all client sessions' history on a file
	clientFileLogger log.Logger
	// Log file that the clientFileLogger operates on.  Stored in this field
	//   so that it can be properly closed when client app is exited.
	clientLogFile    *os.File
}

// Constructor function that creates a production-level logger for logging all client
//   activity, DEBUG-level and above.
func NewClientLoggingObjectDEV(clientLogFilename string) *ClientLoggingObjectDEV{
	// Create terminal logger
	terminalLogger := log.NewLogfmtLogger(os.Stdout)
	terminalLogger = level.NewFilter(terminalLogger, level.AllowDebug())
	terminalLogger = log.With(terminalLogger, "time", log.DefaultTimestampUTC)

	// Add filepath to clientlogs directory to get full client log file path.
	clientLogFilename = "github.com/astronomical3/fewer_grpc/client/clientlogs/devtest/" + clientLogFilename

	// Create or open the client log file object
	clientLogFile, err := os.OpenFile(clientLogFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic (err)
	}

	// Create the logger to operate logging operations on the client log file
	clientFileLogger := log.NewLogfmtLogger(clientLogFile)
	clientFileLogger = level.NewFilter(clientFileLogger, level.AllowDebug())
	clientFileLogger = log.With(clientFileLogger, "time", log.DefaultTimestampUTC)

	return &ClientLoggingObjectDEV{
		terminalLogger: terminalLogger,
		clientFileLogger: clientFileLogger,
		clientLogFile: clientLogFile,
	}
}

// Method of the ClientLoggingObjectDEV that is used for simultaneously logging DEBUG-level
//   activity on both the terminal and client log file.  This is an operation only used when
//   a ClientLoggingObjectDev for a ClientLogger.
func (clod *ClientLoggingObjectDEV) ClientLogDebug(key, value, message string) {
	go level.Debug(clod.terminalLogger).Log(key, value, "message", message)
	level.Debug(clod.clientFileLogger).Log(key, value, "message", message)
}

