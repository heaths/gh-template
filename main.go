// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/auth"
	"github.com/cli/go-gh/pkg/repository"
	"github.com/heaths/gh-template/internal/cmd"
	"github.com/heaths/go-console"
	"github.com/spf13/cobra"
)

func main() {
	con := console.System()

	var repo string
	opts := &cmd.GlobalOptions{
		Console: con,
		Log:     log.New(con.Stdout(), "", log.Ltime),
	}

	rootCmd := &cobra.Command{
		Use:   "template",
		Short: "Format project templates",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if repo != "" {
				opts.Repo, err = repository.Parse(repo)
				if err != nil {
					return
				}
			}

			if opts.Repo == nil {
				opts.Repo, err = gh.CurrentRepository()
				if err != nil {
					return
				}
			}

			if opts.Repo == nil {
				return fmt.Errorf("no repository")
			}

			// Make sure the user is authenticated.
			host := opts.Repo.Host()
			if host == "" {
				host, _ = auth.DefaultHost()
			}

			token, _ := auth.TokenForHost(host)
			if token == "" {
				return fmt.Errorf("use `gh auth login` to authenticate with required scopes")
			}

			return
		},
	}

	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "R", "", "Select another repository to use using the [HOST/]OWNER/REPO format")
	rootCmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false, "Log verbose output")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
