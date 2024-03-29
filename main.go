// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package main

import (
	"log"
	"os"

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
		Long:  "GitHub CLI extension to list and clone repository templates.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if repo != "" {
				opts.Repo, err = repository.Parse(repo)
				if err != nil {
					return
				}
			}

			return
		},
		SilenceUsage: true,
	}

	width, _, err := con.Size()
	if err != nil {
		width = 80
	}
	rootCmd.SetOut(con.Stdout())
	rootCmd.SetHelpFunc(cmd.AppendHelpFunc(width, rootCmd.HelpFunc()))

	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "R", "", "Select another repository to use using the [HOST/]OWNER/REPO format")
	rootCmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false, "Log verbose output")

	rootCmd.AddCommand(cmd.ApplyCmd(opts))
	rootCmd.AddCommand(cmd.CloneCmd(opts))
	rootCmd.AddCommand(cmd.ListCmd(opts))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
