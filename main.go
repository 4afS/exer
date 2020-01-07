package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	runOption   = flag.Bool("run", false, "Run a project you're in.")
	buildOption = flag.Bool("build", false, "Run a project you're in.")
)

func usage() {
	fmt.Printf(`Usage of exer:
  -build
      Run a project you're in.
  -run
      Run a project you're in.
`)
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *runOption && *buildOption {
		fmt.Fprintln(os.Stderr, fail("specify only `run` or `build`"))
		os.Exit(1)
	}

	path, err := projectRootPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, fail(err.Error()))
		os.Exit(1)
	}

	fileinfos, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, fail(err.Error()))
		os.Exit(1)
	}

	command, ok := findCommands(fileinfos)
	if !ok {
		fmt.Fprintln(os.Stderr, fail("this language not supported"))
		os.Exit(1)
	}

	switch {
	case *runOption:
		runCommand, ok := command["run"]
		if !ok {
			fmt.Fprintln(os.Stderr, fail("run command not found"))
			os.Exit(1)
		}

		fmt.Println(success(runCommand))

	case *buildOption:
		buildCommand, ok := command["build"]
		if !ok {
			fmt.Fprintln(os.Stderr, fail("build command not found"))
			os.Exit(1)
		}

		fmt.Println(success(buildCommand))

	default:
		flag.Usage()
	}
}

func fail(message string) string {
	return fmt.Sprint("exer: ", message)
}

func success(message string) string {
	return fmt.Sprint("[Success]", message)
}

func projectRootPath() (string, error) {
	var (
		cmd        = "git"
		cmdOptions = []string{"rev-parse", "--show-toplevel"}
	)

	result, err := exec.Command(cmd, cmdOptions...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(".git directory not found")
	}

	return strings.TrimRight(string(result), "\n"), err
}

func findCommands(fileinfos []os.FileInfo) (map[string]string, bool) {
	commands := map[string]map[string]string{
		"stack.yaml": {"build": "stack build", "run": "stack run"},
		"cargo.toml": {"build": "cargo build", "run": "cargo run"},
		".spago":     {"build": "spago build", "run": "spago run"},
		"elm.json":   {"build": "elm reactor"},
		"build.sbt":  {"build": "sbt build", "run": "sbt run"},
	}

	for _, fileinfo := range fileinfos {
		if command, ok := commands[fileinfo.Name()]; ok {
			return command, true
		}
	}

	return nil, false
}
