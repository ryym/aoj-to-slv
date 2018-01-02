package main

import (
	"io"
	"os/exec"
	"regexp"
)

func NewSlv() Slv {
	return &SlvImpl{
		cmd:              "slv",
		runningPrefixRgx: regexp.MustCompile("running\\s[^\\n]+\\n"),
	}
}

type SlvImpl struct {
	cmd              string
	runningPrefixRgx *regexp.Regexp
}

func (s *SlvImpl) New(dir string) error {
	_, err := exec.Command(s.cmd, "new", dir).CombinedOutput()
	return err
}

func (s *SlvImpl) Run(src string, input string) (string, error) {
	cmd := exec.Command(s.cmd, "run", src)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, input)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return s.runningPrefixRgx.ReplaceAllString(string(out), ""), nil
}

func (s *SlvImpl) Test(src string) ([]byte, error) {
	return exec.Command(s.cmd, "test", src).CombinedOutput()
}
