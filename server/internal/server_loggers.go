package internal

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)



//************************************************************************************************
// Definition of a ServerLogger interface that includes common logging operation methods that will
//   be used across multiple concrete logger types.
type ServerLogger interface {
	ServerLogInfo(key, value, message string)
	ServerLogWarn(key, value, message string)
	ServerLogError(key, value, message string)
	Close()
}



//*************************************************************************************************
// Definition of a production-level server activity logger that logs both on the terminal/stdout
//   and a file.
type ServerLoggingObjectPROD struct {
	// Log only current server session activity on standard output
	terminalLogger   log.Logger
	// Log all server sessions' history on a file
	serverFileLogger log.Logger
	// Log file that the serverFileLogger operates on.  Stored in this field so that
	//   it can be properly closed when client app is exited.
	serverLogFile    *os.File
}

// Constructor function that creates a production-level logger for logging all server activity,
//   INFO-level and above.
func NewServerLoggingObjectPROD(serverLogFilename string) *ServerLoggingObjectPROD {
	// Create terminal logger
	terminalLogger := log.NewLogfmtLogger(os.Stdout)
	terminalLogger = level.NewFilter(terminalLogger, level.AllowInfo())
	terminalLogger = log.With(terminalLogger, "time", log.DefaultTimestampUTC)

	// Add filepath to serverlogs directory to get full server log filepath.
	serverLogFilename = "fewer_grpc/server/serverlogs/production" + serverLogFilename

	// Create or open the server log file object
	serverLogFile, err := os.OpenFile(serverLogFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}

	// Create the logger to operate logging operations on the server log file
	serverFileLogger := log.NewLogfmtLogger(serverLogFile)
	serverFileLogger = level.NewFilter(serverFileLogger, level.AllowInfo())
	serverFileLogger = log.With(serverFileLogger, "time", log.DefaultTimestampUTC)

	return &ServerLoggingObjectPROD{
		terminalLogger:   terminalLogger,
		serverFileLogger: serverFileLogger,
		serverLogFile:    serverLogFile,
	}
}

// Method of the ServerLoggingObjectPROD that is used for simultaneously logging
//   INFO-level activity on both the terminal and server log file.
func (slop *ServerLoggingObjectPROD) ServerLogInfo(key, value, message string) {
	go level.Info(slop.terminalLogger).Log(key, value, "message", message)
	level.Info(slop.serverFileLogger).Log(key, value, "message", message)
}

// Method of the ServerLoggingObjectPROD that is used for simultaneously logging
//   WARN-level activity on both the terminal and server log file.
func (slop *ServerLoggingObjectPROD) ServerLogWarn(key, value, message string) {
	go level.Warn(slop.terminalLogger).Log(key, value, "message", message)
	level.Warn(slop.serverFileLogger).Log(key, value, "message", message)
}

// Method of the ServerLoggingObjectPROD that is used for simultaneouslu logging
//   ERROR-level activity on both the terminal and server log file.
func (slop *ServerLoggingObjectPROD) ServerLogError(key, value, message string) {
	go level.Error(slop.terminalLogger).Log(key, value, "error", message)
	level.Error(slop.serverFileLogger).Log(key, value, "error", message)
}

// Method of the ServerLoggingObjectPROD that is used for closing the server log
//   file properly when the server is about to shutdown.
func (slop *ServerLoggingObjectPROD) Close() {
	slop.serverLogFile.Close()
}



//*************************************************************************************************
// Definition of a development-level server activity logger that logs both on the terminal/stdout
//   and a file.
type ServerLoggingObjectDEV struct {
	// Inherit the ServerLoggingObjectPROD methods, because they also satisfy the 
	// ServerLogger interface.
	ServerLoggingObjectPROD
	// Log only current server session activity on standard output
	terminalLogger   log.Logger
	// Log all server sessions' history on a file
	serverFileLogger log.Logger
	// Log file that the serverFileLogger operates on.  Stored in this field so that
	//   it can be properly closed when client app is exited.
	serverLogFile    *os.File
}

// Constructor function that creates a production-level logger for logging all server activity,
//   DEBUG-level and above.
func NewServerLoggingObjectDEV(serverLogFilename string) *ServerLoggingObjectDEV {
	// Create terminal logger
	terminalLogger := log.NewLogfmtLogger(os.Stdout)
	terminalLogger = level.NewFilter(terminalLogger, level.AllowDebug())
	terminalLogger = log.With(terminalLogger, "time", log.DefaultTimestampUTC)

	// Add filepath to serverlogs directory to get full server log filepath.
	serverLogFilename = "fewer_grpc/server/serverlogs/devtest" + serverLogFilename

	// Create or open the server log file object
	serverLogFile, err := os.OpenFile(serverLogFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic (err)
	}

	// Create the logger to operate logging operations on the server log file
	serverFileLogger := log.NewLogfmtLogger(serverLogFile)
	serverFileLogger = level.NewFilter(serverFileLogger, level.AllowDebug())
	serverFileLogger = log.With(serverFileLogger, "time", log.DefaultTimestampUTC)

	return &ServerLoggingObjectDEV{
		terminalLogger:   terminalLogger,
		serverFileLogger: serverFileLogger,
		serverLogFile:    serverLogFile,
	}
}

// Method of the ServerLoggingObjectDEV that is used for simultaneously logging
//   DEBUG-level activity on both the terminal and server log file.
// This is an operation only used when using a ServerLoggingObjectDEV for a 
//   ServerLogger.
func (slod *ServerLoggingObjectDEV) ServerLogDebug(key, value, message string) {
	go level.Debug(slod.terminalLogger).Log(key, value, "message", message)
	level.Debug(slod.serverFileLogger).Log(key, value, "message", message)
}