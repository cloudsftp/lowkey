package main

import (
	"context"
	"dagger/lowkey/internal/dagger"
)

func (l *Lowkey) Build(
	source *dagger.Directory,
	// +optional
	rusthookSource *dagger.Directory,
) *dagger.File {
	return cachedRustBuilder(source, rusthookSource).
		WithExec([]string{"cargo", "build", "--release"}).
		WithExec([]string{"cp", "target/release/lowkey", "/lowkey"}).
		File("/lowkey")
}

func (l *Lowkey) Test(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	rusthookSource *dagger.Directory,
) (string, error) {
	return cachedRustBuilder(source, rusthookSource).
		WithExec([]string{"cargo", "test"}).
		Stdout(ctx)
}

func (l *Lowkey) Lint(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	rusthookSource *dagger.Directory,
) (string, error) {
	return cachedRustBuilder(source, rusthookSource).
		WithExec([]string{"cargo", "clippy", "--", "-D", "warnings"}).
		Stdout(ctx)
}

func cachedRustBuilder(
	source *dagger.Directory,
	rusthookSource *dagger.Directory,
) *dagger.Container {
	source = source.WithoutDirectory("target")

	return dag.Container().
		From("rust:"+RustVersion).
		WithExec([]string{"rustup", "component", "add", "clippy"}).

		// Source Code
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithDirectory("/rusthook", rusthookSource).

		// Caches
		WithMountedCache("/cache/cargo", dag.CacheVolume("rust-packages")).
		WithEnvVariable("CARGO_HOME", "/cache/cargo").
		WithMountedCache("target", dag.CacheVolume("rust-target"))
}
