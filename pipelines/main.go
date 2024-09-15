package main

import (
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
