package interfaces

// Define a interface para o log.
type Log interface {
	Info(format string, a ...any)
	Error(format string, a ...any)
	Warn(format string, a ...any)
	Debug(format string, a ...any)
}
