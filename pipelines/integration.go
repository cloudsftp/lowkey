package main

import (
	"context"

	"dagger/lowkey/internal/dagger"
)

func (l *Lowkey) TestIntegration(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	devServerExecutable *dagger.File,
) (string, error) {
	lowkeyService := l.BuildLowkeyService(ctx, source)
	localDevService := l.BuildLocalDevService(ctx, source, lowkeyService, devServerExecutable)

	return cachedGoBuilder(source.Directory("integration")).
		WithServiceBinding("lowkey-api", lowkeyService).
		WithServiceBinding("local-dev", localDevService).
		WithExec([]string{"go", "test", "-count=1", "./..."}).
		Stdout(ctx)
}

func (l *Lowkey) BuildLowkeyService(
	ctx context.Context,
	source *dagger.Directory,
) *dagger.Service {
	natsService := l.BuildNatsService(ctx)

	return l.
		BuildBaseImage(ctx, source).
		WithFile(".env", source.File(".env")).
		WithServiceBinding("nats", natsService).
		WithExec([]string{"/bin/server"}).
		AsService()
}

func (l *Lowkey) BuildNatsService(ctx context.Context) *dagger.Service {
	return dag.Container().
		From("nats:latest").
		WithExposedPort(4222).
		WithDefaultArgs([]string{
			"--jetstream", "--name", "main",
		}).
		AsService()
}
