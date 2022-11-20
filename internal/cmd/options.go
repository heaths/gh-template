// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package cmd

import (
	"log"

	"github.com/cli/go-gh/pkg/repository"
	"github.com/heaths/go-console"
)

type GlobalOptions struct {
	Console console.Console
	Log     *log.Logger
	Verbose bool

	Repo repository.Repository
}
