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
	"strings"

	"github.com/olekukonko/tablewriter"
)

// renderSimpleTable is used to render any simple table within the Lacework CLI,
// every command should leverage this function unless there are extra customizations
// required, if so, use instead renderCustomTable(). The benefit of this function
// is the ability to switch/update the look and feel of the human-readable format
// across the entire project
func renderSimpleTable(headers []string, data [][]string) string {
	var (
		tblBldr = &strings.Builder{}
		tbl     = tablewriter.NewWriter(tblBldr)
	)
	tbl.SetHeader(headers)
	tbl.SetRowLine(false)
	tbl.SetBorder(false)
	tbl.SetAutoWrapText(true)
	tbl.SetAlignment(tablewriter.ALIGN_LEFT)
	tbl.SetColumnSeparator(" ")
	tbl.AppendBulk(data)
	tbl.Render()
	return tblBldr.String()
}

type tableOption interface {
	apply(t *tablewriter.Table)
}

type tableFunc func(t *tablewriter.Table)

func (fn tableFunc) apply(t *tablewriter.Table) {
	fn(t)
}

// renderCustomTable should be used on special cases where we need to render a table
// with very specific settings, we normally should use renderSimpleTable() as much
// as possible to have consistency across the CLI
func renderCustomTable(headers []string, data [][]string, opts ...tableOption) string {
	var (
		tblBldr = &strings.Builder{}
		tbl     = tablewriter.NewWriter(tblBldr)
	)

	for _, opt := range opts {
		opt.apply(tbl)
	}

	tbl.SetHeader(headers)
	tbl.AppendBulk(data)
	tbl.Render()

	return tblBldr.String()
}

func renderOneLineCustomTable(title, content string, opts ...tableOption) string {
	return renderCustomTable([]string{title}, [][]string{[]string{content}}, opts...)
}
