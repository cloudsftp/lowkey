package main

import (
	"context"

	"dagger/lowkey/internal/dagger"
)

type Lowkey struct{}

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
