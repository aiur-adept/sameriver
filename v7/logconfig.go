package sameriver

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/TwiN/go-color"

	"go.uber.org/atomic"
)

type PrintfLike func(format string, params ...any)

// produce a function which will act like fmt.Sprintf but be silent or not
// based on a supplied boolean value (below the function definition in this
// file you can find all of them used)
func SubLogFunction(
	moduleName string,
	flag bool,
	wrapper func(s string) string) func(s string, params ...any) {

	prefix := fmt.Sprintf("[%s] ", moduleName)
	return func(format string, params ...any) {
		switch {
		case !flag:
			return
		case len(params) == 0:
			Logger.Printf(wrapper(prefix + format))
		default:
			Logger.Printf(wrapper(fmt.Sprintf(prefix+format, params...)))
		}
	}
}

/* uncork when needed?
var logError = SubLogFunction(
	"ERROR", true,
	func(s string) string { return color.InRed(color.InBold(s)) })
*/

// we have come around to using ERROR LOGGING instead of just always panicing.
// overall i think we should panic way less... Sometimes it's better to just log
// the error in red, but *try* to continue to function. SOmething is broken and
// probably/possibly might lead to a panic, but we keep moving and if it does happen,
// the red bold log is right there, with the greppable prefix [ERROR] :)
var logDSLError = SubLogFunction(
	"ERROR", true,
	func(s string) string {
		stack := make([]byte, 4096)
		length := runtime.Stack(stack, false)
		return color.InRed("[Entity Filter DSL]" + s + "\n" + string(stack[:length]))
	})

var logWarning = SubLogFunction(
	"WARNING", true,
	func(s string) string { return color.InYellow(color.InBold(s)) })

var logWarningRateLimited = func(ms int) PrintfLike {
	var flag atomic.Uint32
	return func(format string, params ...any) {
		if flag.CompareAndSwap(0, 1) {
			logWarning(format, params...)
			go func() {
				time.Sleep(time.Duration(ms) * time.Millisecond)
				flag.CompareAndSwap(1, 0)
			}()
		}
	}
}

var DEBUG_EVENTS = os.Getenv("DEBUG_EVENTS") == "true"
var logEvents = SubLogFunction(
	"Events", DEBUG_EVENTS,
	func(s string) string { return color.InWhiteOverPurple(s) })

var DEBUG_GOAP = os.Getenv("DEBUG_GOAP") == "true"
var logGOAPDebug = SubLogFunction(
	"GOAP", DEBUG_GOAP,
	func(s string) string { return s })

var DEBUG_RUNTIME_LIMITER = os.Getenv("DEBUG_RUNTIME_LIMITER") == "true"
var logRuntimeLimiter = SubLogFunction(
	"RuntimeLimiter", DEBUG_RUNTIME_LIMITER,
	func(s string) string { return s })
