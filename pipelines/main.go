package main

import (
	"context"
	"dagger/lowkey/internal/dagger"
	"fmt"
	"sync"
)

type Lowkey struct {
	// +private
	RegistryConfig *dagger.RegistryConfig
}

const (
	RustVersion = "1.81"
	GoVersion   = "1.23"
)

func (l *Lowkey) BuildAndTestAll(
	ctx context.Context,
	source *dagger.Directory,
) error {
	var wg sync.WaitGroup

	errors := make(chan error)

	wg.Add(1)
	go func() {
		_, err := l.Lint(ctx, source)
		if err != nil {
			errors <- err
			return
		}

		_, err = l.Test(ctx, source)
		if err != nil {
			errors <- err
			return
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		l.Build(source)

		l.BuildImage(ctx, source)

		_, err := l.TestIntegration(ctx, source)
		if err != nil {
			errors <- err
			return
		}

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
		case <-done:
			break loop
		case err := <-errors:
			return err
		}
	}
	// Ende

	return nil
}

func (l *Lowkey) Pipeline(
	ctx context.Context,
	source *dagger.Directory,
	actor string,
	token *dagger.Secret,
	host *dagger.Secret,
	username *dagger.Secret,
	key *dagger.Secret,
) error {
	err := l.BuildAndTestAll(ctx, source)

	_, err = l.PublishImage(ctx, source, actor, token)
	if err != nil {
		return err
	}

	_, err = l.Deploy(ctx, host, username, key)
	if err != nil {
		return err
	}

	return nil
}
