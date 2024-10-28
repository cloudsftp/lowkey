package main

import (
	"context"
	"time"

	"dagger/lowkey/internal/dagger"
)

const LocalDevServerVersion = "1.3.32"

func (l *Lowkey) TestIntegration(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	mittlifeSource *dagger.Directory,
	// +optional
	devServerExecutable *dagger.File,
) (string, error) {
	/*
		lowkeyService := l.BuildLowkeyService(ctx, source)
		localDevService := l.BuildLocalDevService(ctx, source, lowkeyService, devServerExecutable)

		return cachedGoBuilder(source.Directory("integration")).
			WithServiceBinding("lowkey-api", lowkeyService).
			WithServiceBinding("local-dev", localDevService).
			WithExec([]string{"go", "test", "-count=1", "./..."}).
			Stdout(ctx)
	*/
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

	return dag.Container().
		From("ubuntu:latest").
		WithExec([]string{"apt", "update"}).
		WithExec([]string{"apt", "upgrade", "-y"}).
		WithExec([]string{"apt", "install", "curl", "-y"}).
		WithExec([]string{"curl", "http://lowkey:6670/hey"}).
		Stdout(ctx)
}

func (l *Lowkey) IntegrationLowkeyService(
	ctx context.Context,
	source *dagger.Directory,
	localDevService *dagger.Service,
	// +optional
	mittlifeSource *dagger.Directory,
) *dagger.Service {
	natsService := l.BuildNatsService(ctx)

	return l.
		buildBaseImage(ctx, source, mittlifeSource).
		WithFile(".env", getEnvFile(source)).
		WithServiceBinding("nats", natsService).
		WithServiceBinding("local-dev", localDevService).
		WithExec([]string{"/server"}).
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

func (m *Lowkey) IntegrationLocalDevService(
	source *dagger.Directory,
	lowkeyService *dagger.Service,
) *dagger.Service {
	return dag.Container().
		From("mittwald/marketplace-local-dev-server:"+LocalDevServerVersion).
		WithFile(".env", getEnvFile(source)).
		WithServiceBinding("lowkey-api", lowkeyService).
		AsService()
}

func (m *Lowkey) IntegrationDriveTests(
	ctx context.Context,
	source *dagger.Directory,
	lowkeyService *dagger.Service,
	localDevService *dagger.Service,
) (string, error) {
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

		// Services
		WithServiceBinding("lowkey-api", lowkeyService).
		WithServiceBinding("local-dev", localDevService).

		// Execute tests
		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
		WithExec([]string{"go", "test", "-count=1", "./..."}).
		Stdout(ctx)
}

func getEnvFile(source *dagger.Directory) *dagger.File {
	return source.File(".env")
}
