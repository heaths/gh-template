// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package cmd

import (
	"fmt"

	"github.com/heaths/gh-template/internal/git"
	"github.com/heaths/go-template"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
)

func ApplyCmd(globalOpts *GlobalOptions) *cobra.Command {
	var lang string
	opts := &applyOptions{}

	cmd := &cobra.Command{
		Use:         "apply",
		Short:       "Apply project template parameters",
		Long:        "Apply parameters to an already cloned repository template. Any parameters not passed to --param will prompt the user for a value. These may include a default value used if the user does not enter a value.",
		Annotations: annotations(),
		Args:        cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			opts.GlobalOptions = globalOpts

			// Always add .github to avoid processing workflows using ${{...}} expressions.
			if opts.exclusions == nil {
				opts.exclusions = make([]string, 0)
			}
			opts.exclusions = append(opts.exclusions, ".github")

			opts.language, err = language.Parse(lang)
			if err != nil {
				return
			}

			if opts.params == nil {
				opts.params = make(map[string]string)
			}

			return apply(opts)
		},
	}

	cmd.Flags().StringSliceVarP(&opts.exclusions, "exclude", "x", nil, "Paths to exclude using case-insensitive comparisons")
	cmd.Flags().StringVarP(&lang, "language", "l", "en", "Language for some template functions")
	cmd.Flags().StringToStringVarP(&opts.params, "param", "p", nil, "Parameters to apply to project template")

	return cmd
}

type applyOptions struct {
	*GlobalOptions

	exclusions []string
	language   language.Tag
	params     map[string]string
}

func apply(opts *applyOptions) error {
	if name, email, err := git.User(); err == nil {
		fmt.Printf("name = %q, email = %q", name, email)
		opts.params["git.name"] = name
		opts.params["git.email"] = email
	} else {
		fmt.Printf("failed to get config: %v", err)
	}

	if opts.Repo != nil {
		opts.params["github.host"] = opts.Repo.Host()
		opts.params["github.owner"] = opts.Repo.Owner()
		opts.params["github.repo"] = opts.Repo.Name()
	}

	return template.Apply(".", opts.params,
		template.WithExclusions(opts.exclusions),
		template.WithLanguage(opts.language),
		template.WithLogger(opts.Log, opts.Verbose),
	)
}
