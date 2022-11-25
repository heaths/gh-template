// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package cmd

import (
	"github.com/heaths/gh-template/internal/git"
	"github.com/heaths/go-template"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
)

func ApplyCmd(globalOpts *GlobalOptions) *cobra.Command {
	opts := &applyOptions{}

	cmd := &cobra.Command{
		Use:         "apply",
		Short:       "Apply project template parameters",
		Long:        "Apply parameters to an already cloned repository template. Any parameters not passed to --param will prompt the user for a value. These may include a default value used if the user does not enter a value.",
		Annotations: annotations(),
		Args:        cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			opts.GlobalOptions = globalOpts
			return apply(opts)
		},
	}

	applyFlags(cmd, opts)
	return cmd
}

func applyFlags(c *cobra.Command, opts *applyOptions) {
	var lang string

	c.PreRunE = func(cmd *cobra.Command, args []string) (err error) {
		// Always add .github to avoid processing workflows using ${{...}} expressions.
		if opts.exclusions == nil {
			opts.exclusions = make([]string, 0, 1)
		}
		opts.exclusions = append(opts.exclusions, ".github")

		opts.language, err = language.Parse(lang)
		if err != nil {
			return
		}

		if opts.params == nil {
			opts.params = make(map[string]string)
		}

		return
	}

	c.Flags().StringSliceVarP(&opts.exclusions, "exclude", "x", nil, "Any `paths` to exclude using case-insensitive comparisons")
	c.Flags().StringVarP(&lang, "language", "l", "en", "BCP-47 language for some template functions")
	c.Flags().StringToStringVarP(&opts.params, "param", "p", nil, "Parameters to apply to project template as `name=value`")
}

type applyOptions struct {
	*GlobalOptions

	exclusions []string
	language   language.Tag
	params     map[string]string
}

func apply(opts *applyOptions) error {
	if name, email, err := git.User(); err == nil {
		opts.params["git.name"] = name
		opts.params["git.email"] = email
	} else if opts.Verbose && opts.Log != nil {
		opts.Log.Printf("failed to get git config: %v", err)
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
