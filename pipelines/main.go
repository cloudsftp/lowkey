package main

import (
	"context"

	"dagger/lowkey/internal/dagger"
)

type Lowkey struct {
	// +private
	RegistryConfig *dagger.RegistryConfig
}

func (l *Lowkey) Build(source *dagger.Directory) (*dagger.File, error) {
	source = filterDirectory(source)

	builder := cachedRustBuilder(source).
		WithExec([]string{"cargo", "build", "--release"})

	output := builder.File("target/release/lowkey")
	return output, nil
}

func (l *Lowkey) BuildImage(
	ctx context.Context,
	source *dagger.Directory,
) *dagger.Container {
	source = filterDirectory(source)

	container := dag.Container().
		WithDirectory("/src", source).
		WithWorkdir("/src").
		Directory("/src").
		DockerBuild()

	return container
}

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

func filterDirectory(input *dagger.Directory) *dagger.Directory {
	return input.WithoutDirectory("target")
}

func cachedRustBuilder(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("rust:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithMountedCache("/cache/cargo", dag.CacheVolume("rust-cache")).
		WithEnvVariable("CARGO_HOME", "/cache/cargo")
}
