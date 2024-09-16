package main

import (
	"context"
	"dagger/lowkey/internal/dagger"
)

func (l *Lowkey) BuildLocalDevService(
	ctx context.Context,
	source *dagger.Directory,
	lowkeyService *dagger.Service,
	// +optional
	devServerExecutable *dagger.File,
) *dagger.Service {
	if devServerExecutable == nil {
		return buildLocalDevServiceFromImage(
			ctx,
			source,
			lowkeyService,
		)
	} else {
		return buildLocalDevServiceFromSource(
			ctx,
			source,
			devServerExecutable,
			lowkeyService,
		)
	}
}

func buildLocalDevServiceFromImage(
	ctx context.Context,
	source *dagger.Directory,
	lowkeyService *dagger.Service,
) *dagger.Service {
	return dag.
		Container().
		From("mittwald/marketplace-local-dev-server:1.3.6").
		WithFile(".env", source.File(".env")).
		WithServiceBinding("lowkey-api", lowkeyService).
		AsService()
}

func buildLocalDevServiceFromSource(
	ctx context.Context,
	source *dagger.Directory,
	devServerExecutable *dagger.File,
	lowkeyService *dagger.Service,
) *dagger.Service {
	return dag.
		Container().
		From("debian:bookworm-slim").
		WithExposedPort(8080).
		WithFile(".env", source.File(".env")).
		WithFile("local-dev-server", devServerExecutable).
		WithServiceBinding("lowkey-api", lowkeyService).
		WithExec([]string{"./local-dev-server"}).
		AsService()
}
