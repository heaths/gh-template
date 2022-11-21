// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package cmd

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/pkg/repository"
	"github.com/heaths/go-console"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

// cspell:ignore Docf sdescription
func TestList(t *testing.T) {
	tests := []struct {
		name       string
		opts       listOptions
		tty        bool
		mocks      func()
		wantStdout string
		wantErr    bool
	}{
		{
			name: "repositories",
			mocks: func() {
				gock.New("https://api.github.com").
					Post("/graphql").
					Reply(200).
					JSON(`{
						"data": {
							"repositoryOwner": {
								"repositories": {
									"nodes": [
										{
											"owner": {
												"login": "a"
											},
											"name": "c",
											"description": "description c",
											"isTemplate": true
										},
										{
											"owner": {
												"login": "a"
											},
											"name": "z",
											"description": null,
											"isTemplate": false
										}
									],
									"pageInfo": {
										"hasNextPage": true,
										"endCursor": "PAGE_1"
									}
								}
							}
						}
					}`)
				gock.New("https://api.github.com").
					Post("/graphql").
					Reply(200).
					JSON(`{
						"data": {
							"repositoryOwner": {
								"repositories": {
									"nodes": [
										{
											"owner": {
												"login": "A"
											},
											"name": "a",
											"description": "description a",
											"isTemplate": false
										},
										{
											"owner": {
												"login": "a"
											},
											"name": "z",
											"description": null,
											"isTemplate": true
										}
									],
									"pageInfo": {
										"hasNextPage": false,
										"endCursor": null
									}
								}
							}
						}
					}`)
			},
			wantStdout: heredoc.Docf(`
			a/c%[1]sdescription c
			a/z%[1]s
			`, "\t"),
		},
		{
			name: "repositories (TTY)",
			tty:  true,
			mocks: func() {
				gock.New("https://api.github.com").
					Post("/graphql").
					Reply(200).
					JSON(`{
						"data": {
							"repositoryOwner": {
								"repositories": {
									"nodes": [
										{
											"owner": {
												"login": "a"
											},
											"name": "b",
											"description": "description b",
											"isTemplate": true
										},
										{
											"owner": {
												"login": "a"
											},
											"name": "Z",
											"description": "description Z",
											"isTemplate": true
										}
									],
									"pageInfo": {
										"hasNextPage": false,
										"endCursor": null
									}
								}
							}
						}
					}`)
			},
			wantStdout: heredoc.Docf(`
			%[1]s[0;32ma/Z%[1]s[0m  description Z
			%[1]s[0;32ma/b%[1]s[0m  description b
			`, "\033"),
		},
		{
			name: "starred repositories",
			opts: listOptions{
				starred: true,
			},
			mocks: func() {
				gock.New("https://api.github.com").
					Post("/graphql").
					Reply(200).
					JSON(`{
						"data": {
							"repositoryOwner": {
								"repositories": {
									"nodes": [
										{
											"owner": {
												"login": "a"
											},
											"name": "c",
											"description": "description c",
											"isTemplate": true
										},
										{
											"owner": {
												"login": "a"
											},
											"name": "Z",
											"description": "description Z",
											"isTemplate": true
										}
									],
									"pageInfo": {
										"hasNextPage": false,
										"endCursor": null
									}
								}
							}
						}
					}`)
				gock.New("https://api.github.com").
					Post("/graphql").
					Reply(200).
					JSON(`{
						"data": {
							"viewer": {
								"repositories": {
									"nodes": [
										{
											"owner": {
												"login": "b"
											},
											"name": "b",
											"description": null,
											"isTemplate": true
										}
									],
									"pageInfo": {
										"hasNextPage": false,
										"endCursor": null
									}
								}
							}
						}
					}`)
			},
			wantStdout: heredoc.Docf(`
			a/Z%[1]sdescription Z
			a/c%[1]sdescription c
			b/b%[1]s
			`, "\t"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(gock.Off)

			fake := console.Fake(
				console.WithStdoutTTY(tt.tty),
				console.WithSize(120, 40),
			)
			repo, err := repository.Parse("heaths/gh-template")
			assert.NoError(t, err)

			globalOpts := &GlobalOptions{
				Console: fake,
				Repo:    repo,
			}
			tt.opts.GlobalOptions = globalOpts

			if tt.mocks != nil {
				tt.mocks()
			}

			err = list(&tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.True(t, gock.IsDone(), pendingMocks(gock.Pending()))

			stdout, _, _ := fake.Buffers()
			assert.Equal(t, tt.wantStdout, stdout.String())
		})
	}
}

func TestRepositoryNodes_Sort(t *testing.T) {
	t.Parallel()

	sut := repositoryNodes{
		{
			Owner: struct{ Login string }{
				"z",
			},
			Name: "z",
		},
		{
			Owner: struct{ Login string }{
				"Z",
			},
			Name: "z",
		},
		{
			Owner: struct{ Login string }{
				"a",
			},
			Name: "c",
		},
		{
			Owner: struct{ Login string }{
				"a",
			},
			Name:        "b",
			Description: "shouldn't matter",
		},
	}

	sort.Sort(sut)

	repos := make([]string, len(sut))
	for i, s := range sut {
		repos[i] = s.Repo()
	}

	expected := []string{
		"Z/z",
		"a/b",
		"a/c",
		"z/z",
	}

	assert.Equal(t, expected, repos)
}

func pendingMocks(mocks []gock.Mock) string {
	paths := make([]string, len(mocks))
	for i, mock := range mocks {
		paths[i] = mock.Request().URLStruct.String()
	}

	return fmt.Sprintf("%d unmatched mocks: %s", len(paths), strings.Join(paths, ", "))
}
