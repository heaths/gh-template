// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package cmd

import (
	"fmt"
	"log"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/auth"
	"github.com/cli/go-gh/pkg/repository"
	"github.com/heaths/go-console"
)

type GlobalOptions struct {
	Console console.Console
	Log     *log.Logger
	Verbose bool

	Repo repository.Repository

	// Test-only options.
	host      string
	authToken string
}

func (opts *GlobalOptions) EnsureRepository() (err error) {
	if opts.Repo == nil {
		opts.Repo, err = gh.CurrentRepository()
		if err != nil {
			return
		}
	}

	if opts.Repo == nil {
		return fmt.Errorf("no repository")
	}

	return
}

func (opts *GlobalOptions) IsAuthenticated() error {
	// Make sure the user is authenticated.
	host := opts.Repo.Host()
	if host == "" {
		host, _ = auth.DefaultHost()
	}

	token, _ := auth.TokenForHost(host)
	if token == "" {
		return fmt.Errorf("use `gh auth login` to authenticate with required scopes")
	}

	return nil
}
