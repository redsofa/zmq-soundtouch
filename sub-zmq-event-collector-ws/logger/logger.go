package logger

import (
	"io"
	"log"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

/*
 * Source : http://www.goinggo.net/2013/11/using-log-package-in-go.html
 *
 * Usage Example :
 *
 * func main() {
 *   logger.InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
 *	 logger.Info.Printf("Some Info ...")
 *   logger.Error.Println("Something has failed")
 *   //Warning, Trace etc...
 * }
 *
 */
func InitLogger(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
