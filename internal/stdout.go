package internal

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// ReportLogger writes a simple, human readable, line to the provided io.Writer for each
// ReportItem received.
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

func (rl *ReportLogger) Write(_ context.Context, ri ReportItem) error {
	rl.counter++
	_, err := fmt.Fprintf(rl.writer, "%6d - realised %+f on %s\n", rl.counter, ri.RealisedPnL(), ri.SellTimestamp.Format(time.RFC3339))
	return err
}
