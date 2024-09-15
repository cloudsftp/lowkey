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

func (l *Lowkey) Pipeline(
	ctx context.Context,
	source *dagger.Directory,
) error {
	var wg sync.WaitGroup

	errors := make(chan error)

	wg.Add(1)
	go func() {
		result, err := l.Lint(ctx, source)
		if err != nil {
			errors <- err
			return
		}

		fmt.Println(result)

		result, err = l.Test(ctx, source)
		if err != nil {
			errors <- err
			return
		}

		fmt.Println(result)

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		l.Build(source)

		l.BuildImage(ctx, source)

		result, err := l.TestIntegration(ctx, source)
		if err != nil {
			errors <- err
			return
		}

		fmt.Println(result)

		wg.Done()
	}()

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

	return nil
}
