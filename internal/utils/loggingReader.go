package utils

import (
	"fmt"
	"io"
	"time"
)

type LoggingReader struct {
	r        io.ReadCloser
	read     int64
	logger   *Logger
	filename string

	interval time.Duration
	nextLog  time.Time
}

func NewLoggingReader(r io.ReadCloser, logger *Logger, filename string, interval time.Duration) *LoggingReader {
	return &LoggingReader{
		r:        r,
		logger:   logger,
		filename: filename,
		interval: interval,
		nextLog:  time.Now().Add(interval),
	}
}

func (l *LoggingReader) Read(p []byte) (int, error) {
	n, err := l.r.Read(p)
	l.read += int64(n)

	info := Colours["brightBlue"]
	success := Colours["green"]
	reset := Colours["reset"]
	fileName := fmt.Sprintf("`%v%v%v`", Colours["brightYellow"], l.filename, info)

	now := time.Now()
	if l.logger != nil && now.After(l.nextLog) {
		l.logger.Printf("%vcopying %v to file system %v| %v%v", info, fileName, Colours["brightBlack"], FormatBytes(l.read), reset)
		l.nextLog = now.Add(l.interval)
	}

	if err == io.EOF && l.logger != nil {
		l.logger.Printf("%vfinished copying file %v %v| %v%v", success, l.filename, Colours["brightBlack"], FormatBytes(l.read), reset)
	}

	return n, err
}


func (l *LoggingReader) Close() error {
	return l.r.Close()
}
