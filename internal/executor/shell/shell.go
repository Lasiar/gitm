package shell

import (
	"context"
	"os"
	"os/exec"

	"gitm/internal/executor"
)

type Shell struct{}

func New() *Shell {
	return &Shell{}
}

const prefixTempDir = "gitm"

func (s *Shell) Run(ctx context.Context, cmd executor.Command) error {

	c := exec.Command("/bin/bash", "-c", cmd.Script)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout

	return c.Run()
}
