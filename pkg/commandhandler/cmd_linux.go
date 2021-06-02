package commandhandler

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"syscall"
)

func execWait(cmd, workdir string, stdout, stderr io.WriteCloser) (int, error) {
	defer stdout.Close()
	defer stderr.Close()
	return execCommon(workdir, stdout, stderr, func(c *exec.Cmd) error {
		return c.Run()
	}, cmd)
}

func execDontWait(cmd, workdir string) (int, error) {
	// passing '&' as a trailing parameter to /bin/sh in addition (*exec.Command).Start() to will double fork and prevent zombie processes
	return execCommon(workdir, os.Stdout, os.Stderr, func(c *exec.Cmd) error {
		return c.Start()
	}, cmd, "&")
}

func execCommon(workdir string, stdout, stderr io.WriteCloser, execMethodToCall func(*exec.Cmd) error, args ...string) (int, error) {
	args = append([]string{"-c"}, args...)
	c := exec.Command("/bin/sh", args...)
	c.Dir = workdir
	c.Stdout = stdout
	c.Stderr = stderr

	err := execMethodToCall(c)
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			code := status.ExitStatus()
			return code, fmt.Errorf("command terminated with exit status=%d", code)
		}
	}
	if err != nil {
		return 1, errors.Wrapf(err, "failed to execute command")
	}
	return 0, nil
}