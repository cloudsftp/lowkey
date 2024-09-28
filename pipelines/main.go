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
)

// BuildAndTestAll checks, builds, lints, and tests the lowkey service completely
func (l *Lowkey) BuildAndTestAll(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	mittlifeCyclesSource *dagger.Directory,
	// +optional
	devServerExecutable *dagger.File,
) error {
	var wg sync.WaitGroup

	errors := make(chan error)

	wg.Add(1)
	go func() {
		_, err := l.Lint(ctx, source, mittlifeCyclesSource)
		if err != nil {
			errors <- err
			return
		}

		_, err = l.Test(ctx, source, mittlifeCyclesSource)
		if err != nil {
			errors <- err
			return
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		l.Build(source, mittlifeCyclesSource)

		l.BuildImage(ctx, source, mittlifeCyclesSource)

		/*
			_, err := l.TestIntegration(ctx, source, mittlifeCyclesSource, devServerExecutable)
			if err != nil {
				errors <- err
				return
			}
		*/

		wg.Done()
	}()

	// TODO: das muss doch besser gehen :(
	done := make(chan any)
	go func() {
		wg.Wait()
		close(done)
	}()

loop:
	for {
		select {
		case err := <-errors:
			return err
		case <-done:
			break loop
		}
	}
	// Ende

	return nil
}

func (l *Lowkey) PublishAndDeploy(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	mittlifeCyclesSource *dagger.Directory,
	actor string,
	token *dagger.Secret,
	host *dagger.Secret,
	username *dagger.Secret,
	key *dagger.Secret,
) error {
	_, err := l.PublishImage(ctx, source, mittlifeCyclesSource, actor, token)
	if err != nil {
		return err
	}

	_, err = l.Deploy(ctx, host, username, key)
	if err != nil {
		return err
	}

	return nil
}
