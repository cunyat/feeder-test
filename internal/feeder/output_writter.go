package feeder

import (
	"fmt"
	"os"
	"time"
)

// LogWritter is the responsible for writting obtained skus to the given output file
type LogWritter struct {
	file *os.File
}

// NewLogWritter creates a new instance of OutputWritter
func NewLogWritter() (*LogWritter, error) {
	filename := fmt.Sprintf("skus-%d.log", time.Now().Unix())
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not open output file: %w", err)
	}

	return &LogWritter{file: file}, nil
}

// Close closes opened file
func (w *LogWritter) Close() error {
	return w.file.Close()
}

// Write writes data from io.Reader and copies it into the output file
// it will override any previous existing file
func (w *LogWritter) Write(value string) {
	_, err := fmt.Fprintf(w.file, "%s - %s", time.Now().Format(time.RFC3339), value)
	if err != nil {
		fmt.Printf("error writting into log: %s", err)
	}
}
