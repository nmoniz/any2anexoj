package internal

//go:generate go tool mockgen -destination=mocks/mocks_gen.go -package=mocks -typed . RecordReader,Record,ReportWriter
