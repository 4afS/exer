package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		runOption   = flag.Bool("run", false, "Run a project you're in.")
		buildOption = flag.Bool("build", false, "Run a project you're in.")
	)

	flag.Parse()

	path, err := projectRootPath()

	if err != nil {
	}

	switch {
	case *runOption && *buildOption:
		usage()
	case *runOption:
		run(path)
	case *buildOption:
		build(path)
	default:
		usage()
	}
}

func usage() {
	fmt.Printf(`Usage of exer:
  -build
      Run a project you're in.
  -run
      Run a project you're in.
`)
	os.Exit(1)
}

func projectRootPath() (string, error) {
	return "", nil
}

func run(path string) error {
	return nil
}

func build(path string) error {
	return nil
}
