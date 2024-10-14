package main

import (
	"context"
	"sync"

	"dagger/lowkey/internal/dagger"
)

type Lowkey struct {
	// +private
	RegistryConfig *dagger.RegistryConfig
}

const (
	RustVersion = "1.81"
	GoVersion   = "1.23"

	AlpineVersion = "3.20"
)

// BuildAndTestAll checks, builds, lints, and tests the lowkey service completely
func (l *Lowkey) BuildAndTestAll(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	mittlifeSource *dagger.Directory,
	// +optional
	devServerExecutable *dagger.File,
) error {
	var wg sync.WaitGroup

	errors := make(chan error)

	wg.Add(1)
	go func() {
		_, err := l.Lint(ctx, source, mittlifeSource)
		if err != nil {
			errors <- err
			return
		}

		_, err = l.Test(ctx, source, mittlifeSource)
		if err != nil {
			errors <- err
			return
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		l.Build(source, mittlifeSource)

		l.BuildImage(ctx, source, mittlifeSource)

		/*
			_, err := l.TestIntegration(ctx, source, mittlifeCyclesSource, devServerExecutable)
			if err != nil {
				errors <- err
				return
			}
		*/

		wg.Done()
	}()

	done := make(chan any)
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errors:
		return err
	case <-done:
	}

	return nil
}

func (l *Lowkey) PublishAndDeploy(
	ctx context.Context,
	source *dagger.Directory,
	actor string,
	token *dagger.Secret,
	host *dagger.Secret,
	username *dagger.Secret,
	key *dagger.Secret,
	// +optional
	mittlifeSource *dagger.Directory,
) error {
	_, err := l.PublishImage(ctx, source, actor, token, mittlifeSource)
	if err != nil {
		return err
	}

	_, err = l.Deploy(ctx, host, username, key)
	if err != nil {
		return err
	}

	return nil
}
