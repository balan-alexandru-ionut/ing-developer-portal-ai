package level

type ErrorLevel int8

const (
	INFO ErrorLevel = iota
	WARN
	ERROR
	FATAL
)
