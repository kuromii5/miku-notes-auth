package prettylog

import (
	"fmt"
	"log/slog"
	"strconv"
)

const (
	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
)

// Color codes for different log levels
var levelColors = map[slog.Level]int{
	slog.LevelDebug: lightGray,
	slog.LevelInfo:  lightGreen,
	slog.LevelWarn:  lightYellow,
	slog.LevelError: lightRed,
}

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}
