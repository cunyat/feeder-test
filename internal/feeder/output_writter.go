package feeder

import (
	"fmt"
	"io"
	"os"
)

// OutputWritter is the responsible for writting obtained skus to the given output file
type OutputWritter struct {
	outfile string
}

// NewOutputWritter creates a new instance of OutputWritter
func NewOutputWritter(outfile string) *OutputWritter {
	return &OutputWritter{outfile: outfile}
}

// Write writes data from io.Reader and copies it into the output file
// it will override any previous existing file
func (w OutputWritter) Write(input io.Reader) error {
	file, err := os.OpenFile(w.outfile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("could not open output file: %w", err)
	}

	_, err = io.Copy(file, input)
	if err != nil {
		return fmt.Errorf("error writting to output file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("error closing file: %wS", err)
	}

	return nil
}
