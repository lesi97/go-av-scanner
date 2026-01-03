package utils

import (
	"bufio"
	"strings"
)

func ReadLines(sc *bufio.Scanner, logLine func(string)) string {
	var b strings.Builder
	for sc.Scan() {
		line := sc.Text()
		logLine(line)
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}