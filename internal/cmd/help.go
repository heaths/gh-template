// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/cli/go-gh/pkg/tableprinter"
	"github.com/spf13/cobra"
)

const variables = `
git.name	Configured Git user.name
git.email	Configured Git user.email
github.host	The host e.g., github.com
github.owner	Owning user or organization
github.repo	Name of the repository
`

const functions = `
param name [default [prompt]]	Replace name with value, optionally prompting with default
pluralize count thing	Pluralize thing based on count
lowercase string	Make string lowercase
titlecase string	Make string titlecase
uppercase string	Make string uppercase
date	Get UTC date
date.Local	Get local date
date.Year	Get year from date
date.Format layout	Format date based on layout like time.Format()
`

func annotations() map[string]string {
	return map[string]string{
		"help:variables": variables,
		"help:functions": functions,
	}
}

func AppendHelpFunc(original func(*cobra.Command, []string)) func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) {
		original(c, s)

		annotations := c.Annotations
		if annotations != nil {
			if variables, ok := annotations["help:variables"]; ok {
				printAnnotation(c.OutOrStdout(), "Variables:", variables)
			}
			if functions, ok := annotations["help:functions"]; ok {
				printAnnotation(c.OutOrStdout(), "Functions:", functions)
				fmt.Fprintln(c.OutOrStdout())
				fmt.Fprintln(c.OutOrStdout(), "For more information about functions, see https://github.com/heaths/go-template")
			}
		}
	}
}

func printAnnotation(w io.Writer, name, values string) {
	// Print section name.
	fmt.Fprintln(w)
	fmt.Fprintln(w, name)

	// Print value.
	table := tableprinter.New(w, true, 80)
	for _, value := range strings.Split(values, "\n") {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		tokens := strings.SplitN(value, "\t", 2)
		if len(tokens) == 2 {
			table.AddField("  " + tokens[0])
			table.AddField(tokens[1])
			table.EndRow()
		}
	}

	// cspell:ignore nolint errcheck
	table.Render() // nolint:errcheck
}
