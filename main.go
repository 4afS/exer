package main

import (
	"flag"
	"fmt"
	"github.com/mattn/go-shellwords"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
	var (
		runOpt   = flag.Bool("run", false, "Run a project you're in.")
		buildOpt = flag.Bool("build", false, "Run a project you're in.")
	)

	flag.Usage = usage
	flag.Parse()

	if *runOpt && *buildOpt {
		ifFail(fmt.Errorf("select either `run` or `build`"))
	}

	path, err := projectRootPath()
	ifFail(err)

	fileinfos, err := ioutil.ReadDir(path)
	ifFail(err)

	cmd, ok := findCmd(fileinfos)
	if !ok {
		ifFail(fmt.Errorf("this language not supported"))
	}

	err = nil
	switch {
	case *runOpt:
		runCmd, ok := cmd.Run()
		if !ok {
			ifFail(fmt.Errorf("run command not found"))
		}

		err = execute(runCmd, path)

	case *buildOpt:
		buildCmd, ok := cmd.Build()
		if !ok {
			ifFail(fmt.Errorf("build command not found"))
		}

		err = execute(buildCmd, path)

	default:
		flag.Usage()
	}

	ifFail(err)
}

func ifFail(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "exer: %s\n", e.Error())
		os.Exit(1)
	}
}

func execute(cmdstr string, rootPath string) error {
	cmds, err := shellwords.Parse(cmdstr)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd

	switch len(cmds) {
	case 0:
		return fmt.Errorf("unexpected command found")
	case 1:
		cmd = exec.Command(cmds[0])
	default:
		cmd = exec.Command(cmds[0], cmds[1:]...)
	}

	currentPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get path to the current directory")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Chdir(rootPath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("`%s` command not found", cmdstr)
	}
	defer os.Chdir(currentPath)
	return nil
}

func projectRootPath() (string, error) {
	var (
		cmd     = "git"
		cmdOpts = []string{"rev-parse", "--show-toplevel"}
	)

	_, err := exec.LookPath(cmd)
	if err != nil {
		return "", fmt.Errorf("git not installed")
	}

	result, err := exec.Command(cmd, cmdOpts...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(".git directory not found")
	}

	return strings.TrimRight(string(result), "\n"), nil
}

type Cmd struct {
	build, run string
}

func (cmd Cmd) Run() (string, bool) {
	return cmd.run, cmd.run != ""
}

func (cmd Cmd) Build() (string, bool) {
	return cmd.build, cmd.build != ""
}

func findCmd(fileinfos []os.FileInfo) (Cmd, bool) {
	var cmds = map[string]Cmd{
		"stack.yaml":   {build: "stack build", run: "stack run"},
		"Cargo.toml":   {build: "cargo build", run: "cargo run"},
		".spago":       {build: "spago build", run: "spago run"},
		"elm.json":     {run: "elm reactor"},
		"build.sbt":    {build: "sbt build", run: "sbt run"},
		"build.gradle": {build: "gradle build", run: "gradle run"},
	}

	for _, f := range fileinfos {
		if cmd, ok := cmds[f.Name()]; ok {
			return cmd, true
		}
	}

	return Cmd{}, false
}
