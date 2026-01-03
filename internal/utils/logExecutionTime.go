package utils

import (
	"time"
)


func (l *Logger) LogExecutionTime(name string, start time.Time) {
	pathColour := Colours["brightBlack"]
	timeColour := Colours["green"]
	duration := time.Since(start)

	if duration > 100*time.Millisecond {
		timeColour = Colours["brightRed"] + Colours["bold"]
	}

	l.Printf("%s%s %stook %s%s\n", pathColour, name, timeColour, duration, Colours["reset"])
}
