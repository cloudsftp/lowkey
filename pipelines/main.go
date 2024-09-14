package main

import (
	"context"

	"dagger/lowkey/internal/dagger"
)

type Lowkey struct{}

func (l *Lowkey) Build(ctx context.Context, source *dagger.Directory) error {
	rust := dag.Container().From("rust:latest")
	rust = rust.WithDirectory("/work", source).WithWorkdir("/work")

	rust = rust.WithExec([]string{"cargo", "build", "--release"})

	/*
		_, err := rust.AsService().Hostname(ctx)
		if err != nil {
			return err
		}
	*/

	output := rust.File("target/release/lowkey")

	_, err := output.Export(ctx, "result.bin")
	if err != nil {
		return err
	}

	return nil
}
