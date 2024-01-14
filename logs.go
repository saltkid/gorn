package main

import (
	"fmt"
)

const (
	// for informational logs
	INFO = "[INFO] " // no color

	// can safely skip error, doesn't interrupt process
	WARN = "\033[93m[WARN]\033[0m " // yellow

	// cannot safely skip error, must interrupt process
	FATAL = "\033[91m[FATAL]\033[0m " // red

	// for timing purposes
	TIME = "\033[94m[TIME]\033[0m " // blue
)

type LogFlag struct {
	level LogLevel
}
type LogLevel int8 // can only be 1-4
const (
	FATAL_LEVEL LogLevel = iota + 1
	WARN_LEVEL
	INFO_LEVEL
	TIME_LEVEL
)
func (l *LogFlag) Level() (string, error) {
	switch l.level {
	case FATAL_LEVEL:
		return fmt.Sprintln(FATAL), nil
	case WARN_LEVEL:
		return fmt.Sprintln(FATAL, WARN), nil
	case INFO_LEVEL:
		return fmt.Sprintln(FATAL, WARN, INFO), nil
	case TIME_LEVEL:
		return fmt.Sprintln(FATAL, WARN, INFO, TIME), nil
	default:
		return "", fmt.Errorf("invalid log level: %d", l.level)
	}
}
func ToLogLevel(s string) (LogLevel, error) {
	switch s {
	case "fatal":
		return FATAL_LEVEL, nil
	case "warn":
		return WARN_LEVEL, nil
	case "info", "", "all":
		return INFO_LEVEL, nil
	case "time":
		return TIME_LEVEL, nil
	case "none":
		return 0, nil
	}
	return -1, fmt.Errorf("invalid value '%s' for --logs. Must be 'all', 'none', or a valid log header", s)
}
