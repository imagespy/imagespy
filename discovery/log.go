package discovery

// Logger is implemented by a struct that can receive messages.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type noopLogger struct{}

func (n *noopLogger) Debug(args ...interface{})                 {}
func (n *noopLogger) Debugf(format string, args ...interface{}) {}
func (n *noopLogger) Errorf(format string, args ...interface{}) {}

var log Logger = &noopLogger{}

// SetLog allows changing the logger used by other components of the library.
func SetLog(l Logger) {
	log = l
}
