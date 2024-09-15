package main

import (
	"context"

	"dagger/lowkey/internal/dagger"
)

func (l *Lowkey) BuildTestService(
	ctx context.Context,
	source *dagger.Directory,
) *dagger.Service {
	natsService := dag.Container().
		From("nats:latest").
		WithExposedPort(4222).
		WithDefaultArgs([]string{
			"--jetstream", "--name", "main",
		}).
		AsService()

	return l.
		BuildBaseImage(ctx, source).
		WithFile(".env", source.File(".env")).
		WithServiceBinding("nats", natsService).
		WithExec([]string{"/bin/server"}).
		AsService()
}

func buildLocalDevService(
	ctx context.Context,
	source *dagger.Directory,
	lowkeyService *dagger.Service,
) *dagger.Service {
	return dag.
		Container().
		From("mittwald/marketplace-local-dev-server:1.3.5").
		WithFile(".env", source.File(".env")).
		WithServiceBinding("lowkey-api", lowkeyService).
		AsService()
}

func (l *Lowkey) TestIntegration(
	ctx context.Context,
	source *dagger.Directory,
) (string, error) {
	lowkeyService := l.BuildTestService(ctx, source)

	localDevService := buildLocalDevService(ctx, source, lowkeyService)
	localDevService = localDevService

	return cachedGoBuilder(source.Directory("integration")).
		WithServiceBinding("lowkey-api", lowkeyService).
		WithExec([]string{"go", "test"}).
		Stdout(ctx)
}
