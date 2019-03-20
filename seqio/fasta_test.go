package seqio

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

var headerRgxp = regexp.MustCompile(`^\w{1,4}`)
var seqRgxp = regexp.MustCompile(`^[A-Z]+$`)

func TestSingleFasta(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Error("Could not get current direcotry")
	}

	mfasta := filepath.Join(dir, "single.fa")
	reader, err := NewFastaReader(mfasta)
	if err != nil {
		t.Fatalf("error in reading fasta file %s", err)
	}
	if !reader.HasEntry() {
		t.Errorf("Did not get expected iteration %s", err)
	}
	entry := reader.NextEntry()
	if entry == nil {
		t.Error("Did not get expected entry")
	}

	if !bytes.HasPrefix(entry.Id, []byte("tr|Q95Q25")) {
		t.Error("Expected to match header")
	}

	if !seqRgxp.Match(entry.Sequence) {
		t.Error("Expected to match sequence")
	}
	if err := reader.Err(); err != nil {
		t.Errorf("error in reading lines %s", err)
	}
}

func TestMultiFasta(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Error("Could not get current direcotry")
	}

	mfasta := filepath.Join(dir, "multi.fa")
	reader, err := NewFastaReader(mfasta)
	if err != nil {
		t.Fatalf("error in reading fasta file %s", err)
	}
	for i := 0; i <= 3; i++ {
		if !reader.HasEntry() {
			t.Error("Did not get expected iteration")
		}
		entry := reader.NextEntry()
		if entry == nil {
			t.Error("Did not get expected entry")
		}

		if !headerRgxp.Match(entry.Id) {
			t.Error("Expected to match header")
		}

		if !seqRgxp.Match(entry.Sequence) {
			t.Error("Expected to match sequence")
		}
	}
	if err := reader.Err(); err != nil {
		t.Errorf("error in reading lines %s", err)
	}
}
