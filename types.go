package main

import "io"

type ProblemFinder interface {
	ListProblems(dir string) (map[string][]string, error)
}

type Problem interface {
	FindTestFile(files []string) ([]string, string)
	FindFastestSrc(srcs []string) string
}

type TestData interface {
	ListInputs(r io.Reader) ([]string, error) // or Reader
}

type TestItem struct {
	Input  string
	Output string
}

type TestWriter interface {
	Write(w io.Writer, items []*TestItem) error
}

type Slv interface {
	New(dir string) error
	Run(src string, input string) (string, error)
	Test(src string) ([]byte, error)
}

type App struct {
	Slv        Slv
	PrbFinder  ProblemFinder
	TestWriter TestWriter
	TestData   TestData
}
