package dzpk

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kardianos/osext"
	"github.com/mattn/go-isatty"
	"github.com/op/go-logging"
)

var syslog = logging.MustGetLogger("example")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format string

// Password is just an example type implementing the Redactor interface. Any
// time this is logged, the Redacted() function will be called.
// type Password string

// func (p Password) Redacted() interface{} {
// 	return logging.Redact(string(p))
// }

func init() {
	format = "%{time:15:04:05.000} %{shortfunc} > %{level:.4s} %{id:03x} %{message}"

	if (runtime.GOOS == "windows" && len(os.Getenv("ANSICON")) != 0) || isatty.IsTerminal(os.Stdout.Fd()) {
		format = "%{color}%{time:15:04:05.000} %{shortfunc} > %{level:.4s} %{id:03x}%{color:reset} %{message}"
	}
	// For demo purposes, create two backend for os.Stderr.
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	exeFolder, err := osext.ExecutableFolder()
	if err != nil {
		panic(fmt.Sprintf("get executeable folder err:%v", err))
	}

	errFile, err := os.OpenFile(filepath.Join(exeFolder, "err.log"), os.O_APPEND|os.O_CREATE, 0x777)
	if err != nil {
		panic("can't open err log for appending.")
	}
	backend2 := logging.NewLogBackend(errFile, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend1Formatter := logging.NewBackendFormatter(backend1, logging.MustStringFormatter(format))

	// Only errors and more severe messages should be sent to backend1
	backend2Leveled := logging.AddModuleLevel(backend2)
	backend2Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Formatter, backend2Leveled)
}

func GetLogger() *logging.Logger {
	return syslog
}
