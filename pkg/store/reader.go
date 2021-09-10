package store

import "io"

// Reader satisfies io.Reader interface
type Reader struct {
	skus []string
}

var _ io.Reader = &Reader{}

// Read implements io.Reader interface
// this allow to write into file (or others) a large dataset
// readed skus will be deleted from the slice.
func (d *Reader) Read(p []byte) (n int, err error) {
	max := len(p)
	i := 0

	for {
		// if we have read all skus return EOF error
		if i >= len(d.skus) {
			err = io.EOF
			return
		}

		sku := []byte(d.skus[i])

		// check if bytes has space for one more sku
		if len(sku)+n > max {
			break
		}

		// put in output slice sku bytes
		for _, b := range sku {
			p[n] = b
			n++
		}

		i++
	}

	// remove from slice readed skus
	d.skus = append([]string{}, d.skus[i:]...)
	return
}
