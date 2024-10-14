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
	mittlifeSource *dagger.Directory,
) *dagger.File {
	return cachedRustBuilder(source, mittlifeSource).
		WithExec([]string{"cargo", "build", "--release"}).
		WithExec([]string{"cp", "target/release/lowkey", "/lowkey"}).
		File("/lowkey")
}

// Test runs unit tests on the lowkey service
func (l *Lowkey) Test(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	mittlifeSource *dagger.Directory,
) (string, error) {
	return cachedRustBuilder(source, mittlifeSource).
		WithExec([]string{"cargo", "test"}).
		Stdout(ctx)
}

func (l *Lowkey) Lint(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	mittlifeSource *dagger.Directory,
) (string, error) {
	return cachedRustBuilder(source, mittlifeSource).
		WithExec([]string{"cargo", "clippy", "--", "-D", "warnings"}).
		Stdout(ctx)
}

func cachedRustBuilder(
	source *dagger.Directory,
	mittlifeSource *dagger.Directory,
) *dagger.Container {
	source = source.WithoutDirectory("target")

	builder := dag.Container().
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

		// Caches
		WithMountedCache("/cache/cargo", dag.CacheVolume("rust-packages")).
		WithEnvVariable("CARGO_HOME", "/cache/cargo").
		WithMountedCache("target", dag.CacheVolume("rust-target"))

	if mittlifeSource != nil {
		builder = builder.WithDirectory("/mittlife_cycles", mittlifeSource)
	}

	return builder
}
