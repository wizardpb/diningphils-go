package shared

import "fmt"

type Assertion func() bool

// Assert checks an Assertion (as a func() bool) and panics if it is not true
func Assert(f Assertion, msg string, args ...interface{}) {
	if !f() {
		panic(fmt.Sprintf(msg, args))
	}
}
