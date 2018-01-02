package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
)

// e.g.)
// $1 $2 <<-EOF
// 2 2
// 0 0 1
// 1 0 1
// EOF
// echo ------
const TEST_RGX = "\\$1 \\$2 <<-EOF\n((?:[^\\n]*\\n)*?)EOF"

func NewTestData() TestData {
	return &TestDataImpl{
		inputRgx: regexp.MustCompile(TEST_RGX),
	}
}

type TestDataImpl struct {
	inputRgx *regexp.Regexp
}

func (td *TestDataImpl) ListInputs(r io.Reader) ([]string, error) {
	s, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	matchesList := td.inputRgx.FindAllSubmatch(s, -1)

	inputs := make([]string, len(matchesList))
	for i, ms := range matchesList {
		if len(ms) > 0 {
			inputs[i] = string(ms[1])
		}
	}

	return inputs, nil
}

func NewTestWriter() TestWriter {
	return &TestWriterImpl{}
}

type TestWriterImpl struct{}

const QUOTE_FMT = `"""
%s"""
`

func (tw *TestWriterImpl) Write(w io.Writer, items []*TestItem) error {
	var b []byte
	for _, item := range items {
		b = append(b, []byte("[[test]]\nin = ")...)
		b = append(b, []byte(fmt.Sprintf(QUOTE_FMT, item.Input))...)
		b = append(b, []byte("out = ")...)
		b = append(b, []byte(fmt.Sprintf(QUOTE_FMT, item.Output))...)
		b = append(b, []byte("\n")...)
	}
	b = b[:len(b)-1] // Remove last '\n'

	_, err := w.Write(b)
	return err
}
