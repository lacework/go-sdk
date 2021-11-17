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
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra/doc"
)

const fmTemplate = `---
title: "%s"
slug: %s
---

`

func GenerateMarkdownDocs() {
	linkHandler := func(s string) string { return s }
	filePrepender := func(filename string) string {
		var (
			name = filepath.Base(filename)
			base = strings.TrimSuffix(name, path.Ext(name))
		)
		return fmt.Sprintf(fmTemplate, strings.Replace(base, "_", " ", -1), base)
	}

	errcheckEXIT(doc.GenMarkdownTreeCustom(rootCmd, "../docs", filePrepender, linkHandler))
}
