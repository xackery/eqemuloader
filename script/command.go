package script

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func (s *Script) commandRun(command string) (err error) {

	commands := strings.Split(command, " ")
	err = s.commandRunSplit(commands)
	return

}

func (s *Script) commandRunAttached(command string) (err error) {

        commands := strings.Split(command, " ")
        err = s.commandRunSplitAttached(commands)
        return

}

func (s *Script) commandRunSplit(commands []string) (err error) {
	cmd := exec.Command(commands[0], commands[1:len(commands)]...)
	if s.IsVerbose {
		fmt.Println("running command", strings.Join(commands, " "))
	}
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		err = errors.Wrap(err, "failed to start pipe")
		return
	}

	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		err = errors.Wrap(err, "failed to start pipe")
		return
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("> %s\n", scanner.Text())
		}
	}()

	errScanner := bufio.NewScanner(cmdErrReader)
	go func() {
		for errScanner.Scan() {
			fmt.Printf("!> %s\n", errScanner.Text())
		}
	}()
	err = cmd.Start()
	if err != nil {
		err = errors.Wrap(err, "failed to pull image")
		return
	}

	err = cmd.Wait()
	if err != nil {
		err = errors.Wrap(err, "failed during wait")
		return
	}
	return
}

func (s *Script) commandRunParse(command string) (out string, err error) {
	if s.IsVerbose {
		fmt.Println("parsing command", command)
	}
	commands := strings.Split(command, " ")
	cmd := exec.Command(commands[0], commands[1:len(commands)]...)
	var buf bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return
	}
	out = stderr.String()
	if len(out) > 0 {
		return
	}
	out = buf.String()
	return
}

func (s *Script) commandRunDetached(command string) (out string, err error) {
	if s.IsVerbose {
		fmt.Println("running detached command", command)
	}
	commands := strings.Split(command, " ")
	cmd := exec.Command(commands[0], commands[1:len(commands)]...)
	err = cmd.Start()
	if err != nil {
		return
	}
	return
}

func (s *Script) commandRunSplitAttached(commands []string) (err error) {
	cmd := exec.Command(commands[0], commands[1:len(commands)]...)
	if s.IsVerbose {
		fmt.Println("running command", strings.Join(commands, " "))
	}
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		err = errors.Wrap(err, "failed to start pipe")
		return
	}

	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		err = errors.Wrap(err, "failed to start pipe")
		return
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("> %s\n", scanner.Text())
		}
	}()

	errScanner := bufio.NewScanner(cmdErrReader)
	go func() {
		for errScanner.Scan() {
			fmt.Printf("!> %s\n", errScanner.Text())
		}
	}()
	cmd.Stdout = os.Stdout
        cmd.Stdin = os.Stdin
        cmd.Stderr = os.Stderr
        err = cmd.Run()
	if err != nil {
		err = errors.Wrap(err, "failed to pull image")
		return
	}

	err = cmd.Wait()
	if err != nil {
		err = errors.Wrap(err, "failed during wait")
		return
	}
	return
}
