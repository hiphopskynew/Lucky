// It is real package is stderr. We will change it later.
package stdout

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var logging = log.New(os.Stderr, "", 0)

func timestamp() string {
	return time.Now().Format(time.RFC3339)
}

func caller() (string, string, int) {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(4, fpcs)
	if n == 0 {
		fmt.Println("MSG: NO CALLER")
	}
	c := runtime.FuncForPC(fpcs[0] - 1)
	if c == nil {
		fmt.Println("MSG CALLER WAS NIL")
	}
	filepath, line := c.FileLine(fpcs[0] - 1)
	name := c.Name()
	return name, filepath, line
}

func dlog(level string, itf ...interface{}) {
	name, _, line := caller()
	text := []interface{}{}
	text = append(text, fmt.Sprintf("[%v] %v (%v:%v) |", strings.ToUpper(level), timestamp(), name, line))
	text = append(text, itf...)
	logging.Println(text...)
}

func nlog(level string, itf ...interface{}) {
	text := []interface{}{}
	text = append(text, fmt.Sprintf("[%v] %v |", strings.ToUpper(level), timestamp()))
	text = append(text, itf...)
	logging.Println(text...)
}

func Trace(itf ...interface{}) {
	dlog("trace", itf...)
}

func Info(itf ...interface{}) {
	nlog("info", itf...)
}

func Debug(itf ...interface{}) {
	dlog("debug", itf...)
}

func Warning(itf ...interface{}) {
	nlog("warning", itf...)
}

func Error(itf ...interface{}) {
	nlog("error", itf...)
}
