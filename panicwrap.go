// +build !appengine

package bugsnag

import (
	"github.com/kyrylo/bugsnag-go/errors"
	"github.com/mitchellh/panicwrap"
	"os"
)

// NOTE: this function does not return when you call it, instead it
// re-exec()s the current process with panic monitoring.
func defaultPanicHandler() {
	defer defaultNotifier.dontPanic()

	exitStatus, err := panicwrap.Wrap(&panicwrap.WrapConfig{
		CookieKey:   "bugsnag_wrapped",
		CookieValue: "bugsnag_wrapped",
		Handler: func(output string) {

			toNotify, err := errors.ParsePanic(output)

			if err != nil {
				defaultNotifier.Config.log("bugsnag.handleUncaughtPanic: %v", err)
			}
			Notify(toNotify, SeverityError, Configuration{Synchronous: true})
		},
	})

	if err != nil {
		defaultNotifier.Config.log("bugsnag.handleUncaughtPanic: %v", err)
		return
	}

	if exitStatus >= 0 {
		os.Exit(exitStatus)
	} else {
		return
	}
}
