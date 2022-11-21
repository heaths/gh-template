// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/tableprinter"
	"github.com/spf13/cobra"
)

func ListCmd(globalOpts *GlobalOptions) *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists template repositories",
		Long:  "List repository templates from user or organization accounts, and optionally any template repositories the user has starred.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.GlobalOptions = globalOpts

			err := globalOpts.IsAuthenticated()
			if err != nil {
				return err
			}

			return list(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.starred, "starred", "", false, "Include starred repositories")

	return cmd
}

type listOptions struct {
	*GlobalOptions

	starred bool
}

func list(opts *listOptions) (err error) {
	clientOpts := &api.ClientOptions{
		// TODO: Set verbose logging via passthrough buffered writer.
		AuthToken: opts.authToken,
		Host:      opts.host,
	}
	client, err := gh.GQLClient(clientOpts)
	if err != nil {
		return
	}

	templates := repositoryNodes{}
	vars := map[string]interface{}{
		"owner": opts.Repo.Owner(),
	}

	// Loop through associated repos.
	var repos struct {
		RepositoryOwner struct {
			Repositories repositoriesNode
		}
	}
	for {
		err = client.Do(queryRepositories, vars, &repos)
		if err != nil {
			return
		}

		for _, node := range repos.RepositoryOwner.Repositories.Nodes {
			if node.IsTemplate {
				templates = append(templates, node)
			}
		}

		if repos.RepositoryOwner.Repositories.PageInfo.HasNextPage {
			vars["after"] = repos.RepositoryOwner.Repositories.PageInfo.EndCursor
		} else {
			break
		}
	}

	// Loop through starred repos.
	if opts.starred {
		delete(vars, "after")
		var starredRepos struct {
			Viewer struct {
				Repositories repositoriesNode
			}
		}
		for {
			err = client.Do(queryStarredRepositories, vars, &starredRepos)
			if err != nil {
				return
			}

			for _, node := range starredRepos.Viewer.Repositories.Nodes {
				if node.IsTemplate {
					templates = append(templates, node)
				}
			}

			if starredRepos.Viewer.Repositories.PageInfo.HasNextPage {
				vars["after"] = starredRepos.Viewer.Repositories.PageInfo.EndCursor
			} else {
				break
			}
		}
	}

	sort.Sort(templates)

	width := 80
	if opts.Console.IsStdoutTTY() {
		width, _, err = opts.Console.Size()
		if err != nil {
			return
		}
	}

	table := tableprinter.New(opts.Console.Stdout(), opts.Console.IsStdoutTTY(), width)
	cs := opts.Console.ColorScheme()
	for _, template := range templates {
		table.AddField(template.Repo(), tableprinter.WithColor(cs.Green))
		table.AddField(template.Description)
		table.EndRow()
	}
	err = table.Render()

	return
}

type repositoriesNode struct {
	Nodes    repositoryNodes
	PageInfo struct {
		HasNextPage bool
		EndCursor   string
	}
}

type repositoryNode struct {
	Owner struct {
		Login string
	}
	Name        string
	Description string
	IsTemplate  bool

	repo string
}

func (r *repositoryNode) Repo() string {
	if r.repo == "" {
		r.repo = fmt.Sprintf("%s/%s", r.Owner.Login, r.Name)
	}
	return r.repo
}

type repositoryNodes []repositoryNode

func (r repositoryNodes) Len() int {
	return len(r)
}

func (r repositoryNodes) Less(i, j int) bool {
	return strings.Compare(r[i].Repo(), r[j].Repo()) < 0
}

func (r repositoryNodes) Swap(i, j int) {
	r[j], r[i] = r[i], r[j]
}

const queryRepositories = `
query ($owner: String!, $fork: Boolean, $limit: Int = 50, $after: String) {
	repositoryOwner(login: $owner) {
		repositories(isFork: $fork, first: $limit, after: $after) {
			nodes {
				owner {
					login
				}
				name
				description
				isTemplate
			}
			pageInfo {
				hasNextPage
				endCursor
			}
		}
	}
}
`

const queryStarredRepositories = `
query ($limit: Int = 50, $after: String) {
	viewer {
		repositories: starredRepositories(first: $limit, after: $after) {
			nodes {
				owner {
					login
				}
				name
				description
				isTemplate
			}
			pageInfo {
				hasNextPage
				endCursor
			}
		}
	}
}
`
