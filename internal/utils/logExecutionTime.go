package utils

import (
	"fmt"
	"time"
)


func LogExecutionTime(name string, start time.Time) {
	pathColour := Colours["brightBlack"]
	timeColour := Colours["green"]
	duration := time.Since(start)

	if duration > 100*time.Millisecond {
		timeColour = Colours["brightRed"] + Colours["bold"]
	}

	fmt.Printf("%s%s %stook %s%s\n", pathColour, name, timeColour, duration, Colours["reset"])
}
