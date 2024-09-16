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
		BuildBaseImage(ctx, source).
		WithEntrypoint([]string{"/bin/server"})
}

func (l *Lowkey) BuildBaseImage(
	ctx context.Context,
	source *dagger.Directory,
) *dagger.Container {
	executable := l.Build(source)

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
		WithExposedPort(6670).
		WithFile("/bin/server", executable)
}
