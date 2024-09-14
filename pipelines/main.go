package main

import (
	"context"

	"dagger/lowkey/internal/dagger"
)

type Lowkey struct{}

func (l *Lowkey) Check(ctx context.Context, source *dagger.Directory) error {
	rust := dag.Container().From("rust:latest")
	rust = rust.WithDirectory("/work", source).WithWorkdir("/work")

	//path := "."
	//rust = rust.WithExec([]string{"cargo", "check"})
	rust = rust.WithExec([]string{"ls", "-l"})

	_, err := rust.AsService().Hostname(ctx)
	if err != nil {
		return err
	}

	/*
		output := rust.Directory(path)

		_, err = output.Export(ctx, path)
		if err != nil {
			return err
		}
	*/

	return nil
}
