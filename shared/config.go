package shared

import "github.com/frozzare/go-emitter"

// Config represent the
// configurations of the queue
type Config struct {
	Address  string
	Database string
	Emitter  *emitter.Emitter
	Table    string
}
