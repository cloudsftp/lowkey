package main

import (
	"context"

	"dagger/lowkey/internal/dagger"
)

const (
	RustVersion = "1.81"

	APIPort = 6670
)

type Lowkey struct {
	// +private
	RegistryConfig *dagger.RegistryConfig
}

func (l *Lowkey) Build(source *dagger.Directory) *dagger.File {
	return cachedRustBuilder(source).
		WithExec([]string{"cargo", "build", "--release"}).
		WithExec([]string{"cp", "target/release/lowkey", "/lowkey"}).
		File("/lowkey")
}

func (l *Lowkey) Test(
	ctx context.Context,
	source *dagger.Directory,
) (string, error) {
	return cachedRustBuilder(source).
		WithExec([]string{"cargo", "test"}).
		Stdout(ctx)
}

func (l *Lowkey) Lint(
	ctx context.Context,
	source *dagger.Directory,
) (string, error) {
	return cachedRustBuilder(source).
		WithExec([]string{"cargo", "clippy", "--", "-D", "warnings"}).
		Stdout(ctx)
}

func (l *Lowkey) BuildImage(
	ctx context.Context,
	source *dagger.Directory,
) *dagger.Container {
	executable := l.Build(source)
	executablePath := "/bin/server"

	return dag.Container().
		From("debian:bookworm-slim").

		// Install Dependencies
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "libssl3", "ca-certificates"}).
		WithExec([]string{"rm", "-rf", "/var/lib/apt/lists/*"}).

		// User
		WithExec([]string{
			"adduser", "appuser",
			"--disabled-password",
			"--gecos", "",
			"--home", "/nonexistent",
			"--shell", "/sbin/nologin",
			"--no-create-home",
			"--uid", "10001",
		}).
		WithUser("appuser").

		// Application
		WithExposedPort(APIPort).
		WithFile(executablePath, executable).
		WithEntrypoint([]string{executablePath})
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

func cachedRustBuilder(source *dagger.Directory) *dagger.Container {
	source = source.WithoutDirectory("target")

	return dag.Container().
		From("rust:"+RustVersion).
		WithExec([]string{"rustup", "component", "add", "clippy"}).

		// Caches
		WithMountedCache("/cache/cargo", dag.CacheVolume("rust-packages")).
		WithEnvVariable("CARGO_HOME", "/cache/cargo").
		WithMountedCache("target", dag.CacheVolume("rust-target")).

		// Source Code
		WithDirectory("/src", source).
		WithWorkdir("/src")
}

func (l *Lowkey) Deploy(
	ctx context.Context,
	host *dagger.Secret,
	username *dagger.Secret,
	key *dagger.Secret,
) error {
	username_plain, err := username.Plaintext(ctx)
	if err != nil {
		return err
	}

	host_plain, err := host.Plaintext(ctx)
	if err != nil {
		return err
	}

	_, err = dag.SSH().
		Config(username_plain + "@" + host_plain).
		WithIdentityFile(key).
		Command("./deploy.sh").
		AsService().
		Start(ctx)

	if err != nil {
		return err
	}

	return nil
}
