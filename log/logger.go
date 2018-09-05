package log

import (
	"io"
	l "log"
)

// The levels of logging we support
var (
	Trace   *l.Logger
	Info    *l.Logger
	Warning *l.Logger
	Error   *l.Logger
)

// InitLogger : Initialise all the logger function pointers for log levels
func InitLogger(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = l.New(traceHandle,
		"TRACE: ",
		l.Ldate|l.Ltime|l.Lshortfile)

	Info = l.New(infoHandle,
		"INFO: ",
		l.Ldate|l.Ltime|l.Lshortfile)

	Warning = l.New(warningHandle,
		"WARNING: ",
		l.Ldate|l.Ltime|l.Lshortfile)

	Error = l.New(errorHandle,
		"ERROR: ",
		l.Ldate|l.Ltime|l.Lshortfile)
}
