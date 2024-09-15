package main

import (
	"dagger/lowkey/internal/dagger"
)

func cachedGoBuilder(
	source *dagger.Directory,
) *dagger.Container {
	return dag.Container().
		From("golang:"+GoVersion).

		// Caches
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build")).
		WithEnvVariable("GOCACHE", "/go/build-cache").

		// Source code
		WithDirectory("/src", source).
		WithWorkdir("/src")
}
