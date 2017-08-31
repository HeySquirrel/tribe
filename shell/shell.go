package shell

import (
	"errors"
	"fmt"
	tlog "github.com/heysquirrel/tribe/log"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Shell struct {
	pwd    string
	logger *tlog.Log
}

func NewInWd() *Shell {
	pwd, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}

	return New(pwd, tlog.New())
}

func New(pwd string, logger *tlog.Log) *Shell {
	shell := new(Shell)
	shell.pwd = pwd
	shell.logger = logger

	return shell
}

func (shell *Shell) Exec(program string, args ...string) (string, error) {
	shell.logger.Add(fmt.Sprintf("%s %s", program, strings.Join(args, " ")))

	cmd := exec.Command(program, args...)
	cmd.Dir = shell.pwd
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(out), errors.New(err.Error() + " : " + string(out))
	}

	return string(out), err
}
