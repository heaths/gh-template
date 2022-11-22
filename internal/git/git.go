// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

func User() (name, email string, err error) {
	var repo *git.Repository
	if repo, err = git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true}); err != nil {
		return
	}

	return userFromRepo(repo, config.GlobalScope)
}

func userFromRepo(repo *git.Repository, scope config.Scope) (name, email string, err error) {
	var cfg *config.Config
	if cfg, err = repo.ConfigScoped(scope); err != nil {
		return
	}

	if name = cfg.User.Name; name == "" {
		err = fmt.Errorf("user.name not set")
		return
	}
	if email = cfg.User.Email; email == "" {
		err = fmt.Errorf("user.email not set")
		return
	}
	return
}
