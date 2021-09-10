package store

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/cunyat/feeder/pkg/utils"
)

func TestReader_ReadsAllSkus(t *testing.T) {
	r := &StoreReader{
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

	if bytes.Compare(all, expected) != 0 {
		t.Error("given skus and obtained not match")
	}
}

func TestReader_EmptyStore(t *testing.T) {
	r := &StoreReader{}

	all, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("got error reading all: %s", err.Error())
	}

	if len(all) != 0 {
		t.Error("empty store should return emtpy bytes slice")
	}
}
