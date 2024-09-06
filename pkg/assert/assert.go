package assert

import "log"

func Assert(condition bool, msg string) {
	if !condition {
		log.Fatal(msg)
	}
}

func AssertError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
