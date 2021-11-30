//go:generate go run ../docs/main.go
//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (

	// docsLink is the custom link used to render internal links
	docsLink = ""

	// docsCmd is a hidden command that generates automatic documentation in Markdown
	docsCmd = &cobra.Command{
		Use:    "docs <directory>",
		Hidden: true,
		Short:  "Generate Markdown documentation",
		Args:   cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return GenerateMarkdownDocs(args[0])
		},
	}

	// headerTemplate adds front matter to generated documentation, this is how
	// we automatically generate documentation at docs.lacework.com
	headerTemplate = `---
title: "%s"
slug: %s
hide_title: true
---

`
)

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.Flags().StringVarP(&docsLink,
		"link", "l", "", "customize the rendered internal links to the commands")
}

func GenerateMarkdownDocs(location string) error {
	// if the location doesn't exist, we will create it for the user
	if err := os.MkdirAll(location, 0755); err != nil {
		return err
	}

	// given a filename, linkHandler is used to customize the rendered internal links
	// to the commands, only if docsLinks was provided
	linkHandler := func(name string) string {
		if docsLink != "" {
			base := strings.TrimSuffix(name, path.Ext(name))
			return docsLink + strings.ToLower(base) + "/"
		}
		return name
	}

	// filePrepender uses headerTemplate to prepend front matter to the rendered Markdown
	filePrepender := func(filename string) string {
		var (
			name = filepath.Base(filename)
			base = strings.TrimSuffix(name, path.Ext(name))
		)
		return fmt.Sprintf(headerTemplate, strings.Replace(base, "_", " ", -1), base)
	}

	return doc.GenMarkdownTreeCustom(rootCmd, location, filePrepender, linkHandler)
}
