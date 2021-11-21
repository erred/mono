package stdlog

import (
	"errors"
	"log"

	"github.com/go-logr/logr"
)

// New returns a standard library logger,
// logging all messages to l at error level with the given error (or a default message if nil)
func New(l logr.Logger, errmsg error) *log.Logger {
	if errmsg == nil {
		errmsg = errors.New("caught std log")
	}
	return log.New(&writer{errmsg, l}, "", 0)
}

type writer struct {
	errmsg error
	l      logr.Logger
}

func (w *writer) Write(b []byte) (int, error) {
	w.l.Error(w.errmsg, string(b))
	return len(b), nil
}
