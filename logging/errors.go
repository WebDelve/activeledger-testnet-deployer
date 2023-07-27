package logging

import (
	"fmt"

	alsdk "github.com/activeledger/SDK-Golang/v2"
)

func (l *Logger) handleALError(e error, resp alsdk.Response, note string) {
	msg := note
	if len(resp.Summary.Errors) > 0 {
		for i, err := range resp.Summary.Errors {
			msg = fmt.Sprintf("%s\nError %d: %s", msg, i, err)
		}
	}

	msg = fmt.Sprintf("%s\n", msg)

	l.print(msg, ERR, e)
}
