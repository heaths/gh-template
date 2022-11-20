// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package cmd

import (
	"os"

	"github.com/heaths/go-template"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
)

func ApplyCmd(globalOpts *GlobalOptions) *cobra.Command {
	var lang string
	opts := &applyOptions{}

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply project template parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			opts.GlobalOptions = *globalOpts

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

	cmd.Flags().StringVarP(&lang, "language", "l", "en", "Language ")
	cmd.Flags().StringToStringVarP(&opts.params, "param", "p", nil, "Parameters to apply to project template")

	return cmd
}

type applyOptions struct {
	GlobalOptions

	language language.Tag
	params   map[string]string
}

func apply(opts *applyOptions) error {
	opts.params["github.host"] = opts.Repo.Host()
	opts.params["github.owner"] = opts.Repo.Owner()
	opts.params["github.name"] = opts.Repo.Name()

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	return template.Apply(pwd, opts.params,
		template.WithLanguage(opts.language),
	)
}
