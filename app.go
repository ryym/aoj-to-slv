package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const AOJ_PROBLEM_URL = "http://judge.u-aizu.ac.jp/onlinejudge/description.jsp"

func ConvertAll(app *App, dir string) error {
	prbs, err := app.PrbFinder.ListProblems(dir)
	if err != nil {
		return err
	}

	for name, files := range prbs {
		fmt.Printf("%s", name)

		err := Convert(app, dir, name, files)
		if err != nil {
			fmt.Println("")
			return err
		}

		for _, f := range files {
			err = os.Remove(f)
			if err != nil {
				return errors.Wrapf(err, "%s: Failed to remove %s", name, f)
			}
		}

		fmt.Println(" Done")
	}

	return nil
}

func Convert(app *App, dir string, prbName string, prbFiles []string) error {
	prb, err := NewProblem(prbName)
	if err != nil {
		return errors.Wrapf(err, "%s: Failed to create Problem", prbName)
	}

	probDir := filepath.Join(dir, prbName)
	stat, err := os.Stat(probDir)
	if err == nil {
		if stat.IsDir() {
			err = os.RemoveAll(probDir)
			if err != nil {
				return errors.Wrapf(err, "%s: Failed to remove old slv dir", prbName)
			}

			// 既にあるやつは無視する。
			// return nil
		} else {
			return fmt.Errorf("%s: Same name file exists", prbName)
		}
	}

	srcFiles, testFilePath := prb.FindTestFile(prbFiles)
	if len(srcFiles) == 0 {
		return fmt.Errorf("%s :No source files", prbName)
	}

	err = app.Slv.New(probDir)
	if err != nil {
		return errors.Wrapf(err, "%s: Failed to 'slv new'", prbName)
	}

	newSrcFiles := make([]string, len(srcFiles))
	for i, srcFile := range srcFiles {
		srcName := filepath.Base(srcFile)

		f, err := os.Open(srcFile)
		if err != nil {
			return errors.Wrapf(err, "%s: Failed to open %s", prbName, srcName)
		}

		newSrc := filepath.Join(probDir, "src", srcName)
		dest, err := os.Create(newSrc)
		if err != nil {
			return errors.Wrapf(err, "%s: Failed to create %s", prbName, filepath.Join("src", srcName))
		}

		newSrcFiles[i] = newSrc

		_, err = io.Copy(dest, f)
		if err != nil {
			return errors.Wrapf(err, "%s: Failed to copy %s", prbName, srcName)
		}
	}

	if testFilePath != "" {
		testFile, err := os.Open(testFilePath)
		if err != nil {
			return errors.Wrapf(err, "%s: Failed to open test file", prbName)
		}
		inputs, err := app.TestData.ListInputs(testFile)
		if err != nil {
			return errors.Wrapf(err, "%s: Failed to read test inputs", prbName)
		}

		fastSrc := prb.FindFastestSrc(newSrcFiles)
		testItems := make([]*TestItem, len(inputs))
		for i, in := range inputs {
			out, err := app.Slv.Run(fastSrc, in)
			if err != nil {
				return errors.Wrapf(err, "%s: Failed to run with input[%d]: %s", prbName, i, string(out))
			}
			testItems[i] = &TestItem{in, out}
		}

		testDest, err := os.Create(
			filepath.Join(probDir, "test", fmt.Sprintf("%s.toml", prbName)),
		)
		if err != nil {
			return errors.Wrapf(err, "%s: Failed to create test file", prbName)
		}

		err = app.TestWriter.Write(testDest, testItems)
		if err != nil {
			return errors.Wrapf(err, "%s: Failed to write test file", prbName)
		}

		for _, src := range newSrcFiles {
			out, err := app.Slv.Test(src)
			if err != nil {
				return errors.Wrapf(err, "%s: TEST FAIL: %s", prbName, string(out))
			}
		}
	}

	problemId, err := prb.FindProblemId(srcFiles)
	if err != nil {
		return errors.Wrapf(err, "%s: Failed to find problem ID", prbName)
	}
	if problemId != "" {
		ioutil.WriteFile(
			filepath.Join(probDir, "readme.md"),
			[]byte(fmt.Sprintf("<%s?id=%s>\n", AOJ_PROBLEM_URL, problemId)),
			0644,
		)
	}

	return nil
}
