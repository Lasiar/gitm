package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"golang.org/x/sync/errgroup"

	"gitm/internal/config"
	"gitm/internal/executor"
	"gitm/internal/executor/shell"
	shells "gitm/internal/shell"
)

const prefixTempDir = "gitm"

func Run(ctx context.Context, config config.Config) error {
	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(config.ParallelLimit)

	for _, repository := range config.Repositories {
		repository := repository
		eg.Go(func() error {
			return run(egCtx, repository, config.Script)
		})
	}

	return eg.Wait()
}

func run(ctx context.Context, repository string, commands []string) error {
	path, err := os.MkdirTemp("", prefixTempDir)
	if err != nil {
		return err
	}

	_, err = git.PlainCloneContext(ctx, path, false, &git.CloneOptions{
		URL: repository,
	})

	if err != nil {
		return err
	}

	script, err := shells.NewBashShell(path).GenerateScript(commands)
	if err != nil {
		return err
	}
	fmt.Println(script)

	return shell.New().Run(ctx, executor.Command{Script: script})
}
