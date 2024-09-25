package main

import (
	"context"
	"dagger/lowkey/internal/dagger"
)

func (l *Lowkey) Build(
	source *dagger.Directory,
	// +optional
	mittlifeCyclesSource *dagger.Directory,
) *dagger.File {
	return cachedRustBuilder(source, mittlifeCyclesSource).
		WithExec([]string{"cargo", "build", "--release"}).
		WithExec([]string{"cp", "target/release/lowkey", "/lowkey"}).
		File("/lowkey")
}

func (l *Lowkey) Test(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	mittlifeCyclesSource *dagger.Directory,
) (string, error) {
	return cachedRustBuilder(source, mittlifeCyclesSource).
		WithExec([]string{"cargo", "test"}).
		Stdout(ctx)
}

func (l *Lowkey) Lint(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	mittlifeCyclesSource *dagger.Directory,
) (string, error) {
	return cachedRustBuilder(source, mittlifeCyclesSource).
		WithExec([]string{"cargo", "clippy", "--", "-D", "warnings"}).
		Stdout(ctx)
}

func cachedRustBuilder(
	source *dagger.Directory,
	mittlifeCyclesSource *dagger.Directory,
) *dagger.Container {
	source = source.WithoutDirectory("target")

	return dag.Container().
		From("rust:"+RustVersion).
		WithExec([]string{"rustup", "component", "add", "clippy"}).

		// Source Code
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithDirectory("/rusthook", mittlifeCyclesSource).

		// Caches
		WithMountedCache("/cache/cargo", dag.CacheVolume("rust-packages")).
		WithEnvVariable("CARGO_HOME", "/cache/cargo").
		WithMountedCache("target", dag.CacheVolume("rust-target"))
}
