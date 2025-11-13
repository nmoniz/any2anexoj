package internal

type RecordReader interface {
	// ReadRecord should return Records until an error is found.
	ReadRecord() (Record, error)
}

type ReportWriter interface {
	// ReportWriter writes report items
	Write(ReportItem) error
}
