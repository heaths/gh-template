// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package git

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestUserFromRepo(t *testing.T) {
	t.Parallel()

	const (
		_name  = "Test User"
		_email = "test@domain.com"
	)

	type user struct {
		Name  string
		Email string
	}
	tests := []struct {
		name      string
		user      user
		scope     config.Scope
		wantName  bool
		wantEmail bool
		wantErr   bool
	}{
		{
			name: "local",
			user: user{
				Name:  _name,
				Email: _email,
			},
			scope:     config.GlobalScope,
			wantName:  true,
			wantEmail: true,
		},
		{
			name:  "global",
			scope: config.GlobalScope,
		},
		{
			name:    "unset",
			scope:   config.LocalScope,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := memory.NewStorage()
			err := fs.SetConfig(&config.Config{
				User: tt.user,
			})
			assert.NoError(t, err)

			repo, err := git.Init(fs, nil)
			assert.NoError(t, err)

			name, email, err := userFromRepo(repo, tt.scope)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
			}

			if tt.wantName {
				assert.Equal(t, _name, name)
			} else {
				assert.NotEqual(t, _name, name)
			}
			if tt.wantEmail {
				assert.Equal(t, _email, email)
			} else {
				assert.NotEqual(t, _email, email)
			}
		})
	}
}
