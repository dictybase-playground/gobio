//Package seqio is a generic namespace shared by all biological sequence input and output
//handlers.
//This package contain a very barebone and simple Fasta format sequence file parser.
// Currently, it parses and returns the Id(header) and sequence.
// This is mostly working concept, however could be easily extended in future.
// Example
//  package main
//  import (
//		"fmt"
//		"github.com/dictybase-playground/gobio/seqio"
//		"os"
//		"log"
//	)
//
//  func main() {
//		 r,err := seqio.NewFastaReader(os.Args[1])
//		 if err != nil {
//			log.Fatal(err)
//       }
//		 for r.HasEntry()  {
//			fasta := r.NextEentry()
//			fmt.Printf("id:%s\nSequence:%s\n",f.Id,f.Sequence)
//		 }
//       if err := r.Err(); err != nil {
//          fmt.Fprintf(os.Stderr,"error in reading file %s",err)
//       }
//  }
package seqio

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

//A type for holding a single fasta record
type Fasta struct {
	Id       []byte //sequence id or header immediately followd by ">" symbol
	Sequence []byte //The entire sequence
}

//A data type for parsing one entry at a time
type FastaReader struct {
	reader *bufio.Reader //pointer to a buffered reader
	//regular expression for parsing the header. For the time being, it will match any
	//non-whitespace character starting right after the ">" sign in the header.
	seenHeader bool
	header     []byte
	sequence   []byte
	exhausted  bool
	fasta      *Fasta
	err        error
}

//NextEntry the next fasta entry
func (f *FastaReader) NextEntry() *Fasta {
	return f.fasta
}

//Err returns any associated error that happened during the read
func (f *FastaReader) Err() error {
	return f.err
}

//HasEntry checks for next fasta entry.
//Should be called before reading the next entry
func (f *FastaReader) HasEntry() bool {
	for {
		line, err := f.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				if !f.exhausted {
					f.exhausted = true
					f.fasta = &Fasta{Id: f.header, Sequence: f.sequence}
					return true
				}
				return false
			} else {
				f.err = err
				return false
			}
		}
		if bytes.HasPrefix(line, []byte(">")) {
			if !f.seenHeader {
				f.header = line[1 : len(line)-1]
				f.seenHeader = true
			} else {
				f.fasta = &Fasta{Id: f.header, Sequence: f.sequence}
				f.header = line[1 : len(line)-1]
				f.sequence = []byte{}
				return true
			}
		} else {
			f.sequence = append(f.sequence, bytes.TrimSuffix(line, []byte("\n"))...)
		}
	}
	return false
}

//NewFastareader creates a new reader for fasta file
func NewFastaReader(file string) (*FastaReader, error) {
	fr := new(FastaReader)
	reader, err := os.Open(file)
	if err != nil {
		return fr, fmt.Errorf("error in reading %s file", file)
	}
	return &FastaReader{
		reader: bufio.NewReader(reader),
	}, nil
}
