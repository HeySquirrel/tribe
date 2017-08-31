package shell

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

type Shell struct {
	pwd string
}

func NewInWd() *Shell {
	pwd, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}

	return New(pwd)
}

func New(pwd string) *Shell {
	shell := new(Shell)
	shell.pwd = pwd

	return shell
}

func (shell *Shell) Exec(program string, args ...string) (string, error) {
	cmd := exec.Command(program, args...)
	cmd.Dir = shell.pwd
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), errors.New(err.Error() + " : " + string(out))
	}

	return string(out), err
}
