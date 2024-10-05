package main

import (
	"context"
	"fmt"

	"dagger/lowkey/internal/dagger"
)

// Build builds the lowkey service and returns the executable
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

// Test runs unit tests on the lowkey service
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
		From(fmt.Sprintf("rust:%s-alpine%s", RustVersion, AlpineVersion)).
		WithExec([]string{"apk", "update"}).
		WithExec([]string{
			"apk", "add", "--no-cache",
			"pkgconfig", "musl-dev",
			"openssl-dev", "openssl-libs-static",
		}).
		WithExec([]string{"rustup", "component", "add", "clippy"}).

		// Source Code
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithDirectory("/mittlife_cycles", mittlifeCyclesSource).

		// Caches
		WithMountedCache("/cache/cargo", dag.CacheVolume("rust-packages")).
		WithEnvVariable("CARGO_HOME", "/cache/cargo").
		WithMountedCache("target", dag.CacheVolume("rust-target"))
}
