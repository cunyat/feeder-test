package store

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/cunyat/feeder/pkg/utils"
)

func TestReader_ReadsAllSkus(t *testing.T) {
	r := &Reader{
		skus: utils.GenerateSKUs(20),
	}

	all, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("got error reading all: %s", err.Error())
	}

	expected := []byte(strings.Join(r.skus, ""))

	if len(all) != len(expected) {
		t.Error("given skus and obtained bytes not match in length")
	}

	if !bytes.Equal(all, expected) {
		t.Error("given skus and obtained not match")
	}
}

func TestReader_EmptyStore(t *testing.T) {
	r := &Reader{}

	all, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("got error reading all: %s", err.Error())
	}

	if len(all) != 0 {
		t.Error("empty store should return emtpy bytes slice")
	}
}

func BenchmarkReader(b *testing.B) {
	r := &Reader{skus: utils.GenerateSKUs(100000)}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		r := bufio.NewReader(r)
		_, err := r.ReadString('\n')
		for err == nil {
			_, err = r.ReadString('\n')
		}
		if err != io.EOF {
			panic(err)
		}
	}
}
