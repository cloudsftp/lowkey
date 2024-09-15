package main

import (
	"context"

	"dagger/lowkey/internal/dagger"
)

func (l *Lowkey) Deploy(
	ctx context.Context,
	host *dagger.Secret,
	username *dagger.Secret,
	key *dagger.Secret,
) (string, error) {
	usernamePlain, err := username.Plaintext(ctx)
	if err != nil {
		return "", err
	}

	hostPlain, err := host.Plaintext(ctx)
	if err != nil {
		return "", err
	}

	return NewSSH(
		usernamePlain+"@"+hostPlain,
		key,
	).
		Command("./deploy.sh").
		Stdout(ctx)
}
