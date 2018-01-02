package main

import (
	"fmt"
	"os"
)

const BASE_DIR = "/Users/ryu/ghq/github.com/ryym/aoj/"

func main() {
	app := &App{
		Slv:        NewSlv(),
		TestData:   NewTestData(),
		TestWriter: NewTestWriter(),
		PrbFinder:  NewProblemFinder(),
	}

	if len(os.Args) != 2 {
		fmt.Println("Specify sub directory name")
		os.Exit(1)
	}

	subDir := os.Args[1]

	err := ConvertAll(app, BASE_DIR+subDir)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
