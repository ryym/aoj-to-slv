package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func NewProblem(probName string) (Problem, error) {
	return &ProblemImpl{probName}, nil
}

type ProblemImpl struct {
	Name string
}

func (p *ProblemImpl) FindTestFile(files []string) ([]string, string) {
	// No test file
	if len(files) == 1 {
		return files, ""
	}

	srcFiles := make([]string, len(files)-1)
	var testFile string
	i := 0

	for _, f := range files {
		if strings.HasSuffix(f, ".t.sh") {
			testFile = f
		} else {
			srcFiles[i] = f
			i += 1
		}
	}

	return srcFiles, testFile
}

func (p *ProblemImpl) FindFastestSrc(srcs []string) string {
	if len(srcs) == 1 {
		return srcs[0]
	}

	for _, ext := range []string{".cpp", ".rb", ".scala"} {
		for _, src := range srcs {
			if filepath.Ext(src) == ext {
				return src
			}
		}
	}

	return srcs[0]
}

func NewProblemFinder() ProblemFinder {
	return &ProblemFinderImpl{}
}

type ProblemFinderImpl struct{}

func (pf *ProblemFinderImpl) ListProblems(dir string) (map[string][]string, error) {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	prbFiles := make(map[string][]string)
	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		name := removeExts(f.Name())
		fullPath := filepath.Join(dir, f.Name())
		files, exist := prbFiles[name]
		if !exist {
			prbFiles[name] = []string{fullPath}
		} else {
			prbFiles[name] = append(files, fullPath)
		}
	}

	return prbFiles, nil
}

func removeExts(s string) string {
	for true {
		ext := filepath.Ext(s)
		if ext == "" {
			break
		}
		s = strings.TrimSuffix(s, ext)
	}
	return s
}
