// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package cmd

import (
	"fmt"
	"os"

	"github.com/cli/go-gh"
	"github.com/spf13/cobra"
)

func CloneCmd(globalOpts *GlobalOptions) *cobra.Command {
	opts := &cloneOptions{}
	cmd := &cobra.Command{
		Use:         "clone name --template repository",
		Short:       "Clones and formats a template repository",
		Long:        "Clones a template repository then formats any templates found. Any parameters not passed to --param will prompt the user for a value. These may include a default value used if the user does not enter a value.",
		Annotations: annotations(),
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			opts.GlobalOptions = globalOpts

			if len(args) != 1 {
				return fmt.Errorf("expected template repository name")
			}
			opts.name = args[0]

			return clone(opts)
		},
	}

	// Add `apply` flags and parsing, validation pre-run.
	applyFlags(cmd, &opts.applyOptions)

	cmd.Flags().StringVarP(&opts.description, "description", "d", "", "Description of the repository")
	cmd.Flags().StringVar(&opts.template, "template", "", "Make the new `repository` based on a template repository")
	cmd.Flags().StringVarP(&opts.remote, "remote", "r", "", "Specify remote name for the new repository")
	cmd.Flags().StringVar(&opts.homepage, "homepage", "", "Repository home page `URL`")
	cmd.MarkFlagRequired("template") // nolint:errcheck

	cmd.Flags().BoolVar(&opts.disableIssues, "disable-issues", false, "Disable issues in the new repository")
	cmd.Flags().BoolVar(&opts.disableWiki, "disable-wiki", false, "Disable wiki in the new repository")
	cmd.Flags().BoolVar(&opts.includeAllBranches, "include-all-branches", false, "Include all branches from template repository")

	// Do not determine which visibility to use; let `gh repo create` handle that downstream.
	cmd.Flags().BoolVar(&opts.internal, "internal", false, "Make the new repository internal")
	cmd.Flags().BoolVar(&opts.private, "private", false, "Make the new repository private")
	cmd.Flags().BoolVar(&opts.public, "public", false, "Make the new repository public")
	cmd.Flags().StringVarP(&opts.team, "team", "t", "", "The `name` of the organization team to be granted access")
	cmd.MarkFlagsMutuallyExclusive("internal", "private", "public")

	return cmd
}

type cloneOptions struct {
	applyOptions

	name        string
	description string
	template    string
	remote      string
	homepage    string

	disableIssues      bool
	disableWiki        bool
	includeAllBranches bool

	internal bool
	private  bool
	public   bool
	team     string
}

func clone(opts *cloneOptions) (err error) {
	args := make([]string, 0, 18)
	args = append(args, "repo", "create", opts.name, "--template", opts.template, "--clone")
	if opts.description != "" {
		args = append(args, "--description", opts.description)
	}
	if opts.remote != "" {
		args = append(args, "--remote", opts.remote)
	}
	if opts.homepage != "" {
		args = append(args, "--homepage", opts.homepage)
	}
	if opts.disableIssues {
		args = append(args, "--disable-issues")
	}
	if opts.disableWiki {
		args = append(args, "--disable-wiki")
	}
	if opts.includeAllBranches {
		args = append(args, "--include-all-branches")
	}
	if opts.internal {
		args = append(args, "--internal")
	} else if opts.private {
		args = append(args, "--private")
	} else if opts.public {
		args = append(args, "--public")
	}
	if opts.team != "" {
		args = append(args, "--team", opts.team)
	}

	opts.Console.StartProgress("Creating repository " + opts.name)
	_, stderr, err := gh.Exec(args...)
	opts.Console.StopProgress()

	if err != nil {
		fmt.Fprintln(opts.Console.Stderr(), stderr.String())
		return fmt.Errorf("failed to create repository %s: %w", opts.name, err)
	}

	err = os.Chdir(opts.name)
	if err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", opts.name, err)
	}

	// Now that we're in a repo...
	opts.Repo, err = gh.CurrentRepository()
	if err != nil {
		return fmt.Errorf("failed to get repository information: %w", err)
	}

	return apply(&opts.applyOptions)
}
