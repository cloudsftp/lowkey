package main

import (
	"context"

	"dagger/lowkey/internal/dagger"
)

func (l *Lowkey) PublishImage(
	ctx context.Context,
	source *dagger.Directory,
	actor string,
	token *dagger.Secret,
) (string, error) {
	return l.
		BuildImage(ctx, source).
		WithRegistryAuth("ghcr.io", actor, token).
		Publish(ctx, "ghcr.io/cloudsftp/lowkey:latest")
}

func (l *Lowkey) BuildImage(
	ctx context.Context,
	source *dagger.Directory,
) *dagger.Container {
	return l.
		buildBaseImage(ctx, source).
		WithEntrypoint([]string{"/server"})
}

func (l *Lowkey) buildBaseImage(
	ctx context.Context,
	source *dagger.Directory,
) *dagger.Container {
	executable := l.Build(source)

	return dag.Container().
		From("alpine:"+AlpineVersion).
		WithExposedPort(6670).
		WithFile("/server", executable)
}
