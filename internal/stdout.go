package internal

import (
	"fmt"
	"io"
	"os"
	"time"
)

type ReportLogger struct {
	counter int
	writer  io.Writer
}

func NewStdOutLogger() *ReportLogger {
	return &ReportLogger{
		writer: os.Stdout,
	}
}

func NewReportLogger(w io.Writer) *ReportLogger {
	return &ReportLogger{
		writer: w,
	}
}

func (rl *ReportLogger) Write(ri ReportItem) error {
	rl.counter++
	_, err := fmt.Fprintf(rl.writer, "%6d - realised %+f on %s\n", rl.counter, ri.RealisedPnL(), ri.SellTimestamp.Format(time.RFC3339))
	return err
}
