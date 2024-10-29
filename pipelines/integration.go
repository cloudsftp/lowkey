package main

import (
	"context"
	"time"

	"dagger/lowkey/internal/dagger"
)

const LocalDevServerVersion = "latest"

func (l *Lowkey) TestIntegration(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	mittlifeSource *dagger.Directory,
) (string, error) {
	natsService := l.BuildNatsService(ctx)

	_, err := l.buildBaseImage(ctx, source, mittlifeSource).
		WithFile(".env", getEnvFile(source)).
		WithServiceBinding("nats", natsService).
		WithExec([]string{"/server"}).
		AsService().
		WithHostname("lowkey").
		Start(ctx)
	if err != nil {
		return "", err
	}

	_, err = dag.Container().
		From("mittwald/marketplace-local-dev-server:"+LocalDevServerVersion).
		WithFile(".env", getEnvFile(source)).
		AsService().
		WithHostname("local-dev").
		Start(ctx)
	if err != nil {
		return "", err
	}

	return dag.Container().
		From("golang:"+GoVersion).

		// Caches
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build")).
		WithEnvVariable("GOCACHE", "/go/build-cache").

		// Sources
		WithDirectory("/src", source.Directory("integration")).
		WithWorkdir("/src").
		WithFile(".env", getEnvFile(source)).

		// Execute tests
		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
		WithExec([]string{"go", "test", "-count=1", "./..."}).
		Stdout(ctx)
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

func getEnvFile(source *dagger.Directory) *dagger.File {
	return source.File(".env")
}
