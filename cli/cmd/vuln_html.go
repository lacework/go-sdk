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
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/lacework/go-sdk/internal/databox"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/api"
)

const (
	magicBarChartAdjustment    float32 = 5
	magicTableInitialHeight    int32   = 130 // width 1524px can we make it magic?
	magicTableHeightMultiplier int32   = 40
)

// arrays can't be constants
var magicLayerInstructions = []string{"ADD"}

// struct used to render the HTML assessment
type vulnImageAssessmentHtml struct {
	Account     string
	ID          string
	Digest      string
	Repository  string
	CreatedTime string
	Size        string
	Tags        []string

	TotalVulnerabilities    int32
	CriticalVulnerabilities int32
	HighVulnerabilities     int32
	MediumVulnerabilities   int32
	LowVulnerabilities      int32
	InfoVulnerabilities     int32
	FixableVulnerabilities  int32

	TableHeight     int32
	Vulnerabilities []htmlVuln
}

type htmlVuln struct {
	RowHeight         int32
	CVE               string
	Severity          string
	SeverityHTMLClass string
	Layer             string
	PkgName           string
	PkgVersion        string
	PkgFixed          string
	V3Score           float64
	UseV3Score        bool
	V2Score           float64
	UseV2Score        bool
	UseNoScore        bool
}

func (a *vulnImageAssessmentHtml) CriticalVulnPercent() float32 {
	if a.CriticalVulnerabilities == 0 {
		return magicBarChartAdjustment
	}
	percent := (float32(a.CriticalVulnerabilities) / float32(a.TotalVulnerabilities)) * 100
	if a.SeverityToBeAdjusted() == "critical" {
		return percent - a.Adjustment()
	}
	return percent
}
func (a *vulnImageAssessmentHtml) HighVulnPercent() float32 {
	if a.HighVulnerabilities == 0 {
		return magicBarChartAdjustment
	}
	percent := (float32(a.HighVulnerabilities) / float32(a.TotalVulnerabilities)) * 100
	if a.SeverityToBeAdjusted() == "high" {
		return percent - a.Adjustment()
	}
	return percent
}
func (a *vulnImageAssessmentHtml) MediumVulnPercent() float32 {
	if a.MediumVulnerabilities == 0 {
		return magicBarChartAdjustment
	}
	percent := (float32(a.MediumVulnerabilities) / float32(a.TotalVulnerabilities)) * 100
	if a.SeverityToBeAdjusted() == "medium" {
		return percent - a.Adjustment()
	}
	return percent
}
func (a *vulnImageAssessmentHtml) LowVulnPercent() float32 {
	if a.LowVulnerabilities == 0 {
		return magicBarChartAdjustment
	}
	percent := (float32(a.LowVulnerabilities) / float32(a.TotalVulnerabilities)) * 100
	if a.SeverityToBeAdjusted() == "low" {
		return percent - a.Adjustment()
	}
	return percent
}
func (a *vulnImageAssessmentHtml) InfoVulnPercent() float32 {
	if a.InfoVulnerabilities == 0 {
		return magicBarChartAdjustment
	}
	percent := (float32(a.InfoVulnerabilities) / float32(a.TotalVulnerabilities)) * 100
	if a.SeverityToBeAdjusted() == "info" {
		return percent - a.Adjustment()
	}
	return percent
}
func (a *vulnImageAssessmentHtml) Adjustment() float32 {
	var x float32
	if a.CriticalVulnerabilities == 0 {
		x += magicBarChartAdjustment
	}
	if a.HighVulnerabilities == 0 {
		x += magicBarChartAdjustment
	}
	if a.MediumVulnerabilities == 0 {
		x += magicBarChartAdjustment
	}
	if a.LowVulnerabilities == 0 {
		x += magicBarChartAdjustment
	}
	if a.InfoVulnerabilities == 0 {
		x += magicBarChartAdjustment
	}
	return x
}
func (a *vulnImageAssessmentHtml) SeverityToBeAdjusted() string {
	severity := "critical"
	highest := a.CriticalVulnerabilities

	if highest < a.HighVulnerabilities {
		severity = "high"
		highest = a.HighVulnerabilities
	}
	if highest < a.MediumVulnerabilities {
		severity = "medium"
		highest = a.MediumVulnerabilities
	}
	if highest < a.LowVulnerabilities {
		severity = "low"
		highest = a.LowVulnerabilities
	}
	if highest < a.InfoVulnerabilities {
		severity = "info"
	}

	return severity
}

func calcHtmlBarChartWidth(severity string, htmlData vulnImageAssessmentHtml) string {
	switch severity {
	case "critical":
		return fmt.Sprintf("%f%%", htmlData.CriticalVulnPercent())
	case "high":
		return fmt.Sprintf("%f%%", htmlData.HighVulnPercent())
	case "medium":
		return fmt.Sprintf("%f%%", htmlData.MediumVulnPercent())
	case "low":
		return fmt.Sprintf("%f%%", htmlData.LowVulnPercent())
	case "info":
		return fmt.Sprintf("%f%%", htmlData.InfoVulnPercent())
	default:
		return fmt.Sprintf("0")
	}
}

func calcHtmlBarChartX(severity string, htmlData vulnImageAssessmentHtml) string {
	var x float32
	switch severity {
	case "critical":
		x = 0
	case "high":
		x = htmlData.CriticalVulnPercent()
	case "medium":
		x = htmlData.CriticalVulnPercent() + htmlData.HighVulnPercent()
	case "low":
		x = htmlData.CriticalVulnPercent() + htmlData.HighVulnPercent() +
			htmlData.MediumVulnPercent()
	case "info":
		x = htmlData.CriticalVulnPercent() + htmlData.HighVulnPercent() +
			htmlData.MediumVulnPercent() + htmlData.LowVulnPercent()
	}
	return fmt.Sprintf("%f%%", x)
}

func isLayerInstruction(inst string) bool {
	for _, mInst := range magicLayerInstructions {
		if mInst == inst {
			return true
		}
	}
	return false
}
func htmlLayerInstruction(layer string) string {
	if len(layer) == 0 {
		return ""
	}

	words := strings.Split(layer, " ")
	if isLayerInstruction(words[0]) {
		return words[0]
	}

	return "RUN"
}

func htmlLayerPrint(layer string) string {
	if len(layer) == 0 {
		return ""
	}

	words := strings.Split(layer, " ")
	if isLayerInstruction(words[0]) {
		return strings.Join(words[1:], " ")
	}

	return layer
}

func calcHtmlBarChartTextX(severity string, htmlData vulnImageAssessmentHtml) string {
	var x float32
	switch severity {
	case "critical":
		x = htmlData.CriticalVulnPercent() / 2
	case "high":
		x = htmlData.CriticalVulnPercent() + (htmlData.HighVulnPercent() / 2)
	case "medium":
		x = htmlData.CriticalVulnPercent() + htmlData.HighVulnPercent() +
			(htmlData.MediumVulnPercent() / 2)
	case "low":
		x = htmlData.CriticalVulnPercent() + htmlData.HighVulnPercent() +
			htmlData.MediumVulnPercent() + (htmlData.LowVulnPercent() / 2)
	case "info":
		x = htmlData.CriticalVulnPercent() + htmlData.HighVulnPercent() +
			htmlData.MediumVulnPercent() + htmlData.LowVulnPercent() + (htmlData.InfoVulnPercent() / 2)
	}
	return fmt.Sprintf("%f%%", x)
}

func generateVulnAssessmentHTML(assessment *api.VulnContainerAssessment) error {
	// @afiune today, the backend doesn't return any information about the container
	// image if there are NO vulnerabilities, therefore, we are not able to render
	// an HTML file without this data (we skip it)
	if assessment.Image == nil {
		cli.Log.Infow("unable to render assessment in HTML format", "assessment", assessment)
		return nil
	}

	htmlTemplate, ok := databox.Get("vuln_assessment.html")
	if !ok {
		return errors.New(
			"html template not found, this is most likely a mistake on us, please report it to support.lacework.com.",
		)
	}

	var (
		buff    = &bytes.Buffer{}
		funcMap = template.FuncMap{
			"calcBarChartWidth": calcHtmlBarChartWidth,
			"calcBarChartX":     calcHtmlBarChartX,
			"calcBarChartTextX": calcHtmlBarChartTextX,
			"layerInstruction":  htmlLayerInstruction,
			"layerPrint":        htmlLayerPrint,
		}
		tmpl     = template.Must(template.New("vuln_assessment").Funcs(funcMap).Parse(string(htmlTemplate)))
		htmlData = vulnImageAssessmentHtml{
			Account:     cli.Account,
			Repository:  assessment.Image.ImageInfo.Repository,
			ID:          assessment.Image.ImageInfo.ImageID,
			Digest:      assessment.Image.ImageInfo.ImageDigest,
			CreatedTime: assessment.Image.ImageInfo.CreatedTime,
			Size:        byteCountBinary(assessment.Image.ImageInfo.Size),
			Tags:        assessment.Image.ImageInfo.Tags,

			TotalVulnerabilities:    assessment.TotalVulnerabilities,
			CriticalVulnerabilities: assessment.CriticalVulnerabilities,
			HighVulnerabilities:     assessment.HighVulnerabilities,
			MediumVulnerabilities:   assessment.MediumVulnerabilities,
			LowVulnerabilities:      assessment.LowVulnerabilities,
			InfoVulnerabilities:     assessment.InfoVulnerabilities,
			FixableVulnerabilities:  assessment.FixableVulnerabilities,

			TableHeight:     magicTableInitialHeight + (assessment.TotalVulnerabilities * magicTableHeightMultiplier),
			Vulnerabilities: vulContainerImageLayersToHTML(assessment.Image),
		}
		outputHTML = fmt.Sprintf("%s-%s.html",
			strings.ReplaceAll(htmlData.Repository, "/", "-"),
			htmlData.Digest)
	)

	if err := tmpl.Execute(buff, htmlData); err != nil {
		return errors.Wrap(err, "unable to execute template")
	}

	if err := ioutil.WriteFile(outputHTML, buff.Bytes(), os.ModePerm); err != nil {
		return errors.Wrap(err, "unable to write html file")
	}

	cli.OutputHuman("The container vulnerability assessment was stored at '%s'\n", outputHTML)
	return nil
}

func vulContainerImageLayersToHTML(image *api.VulnContainerImage) []htmlVuln {
	if image == nil {
		return []htmlVuln{}
	}

	var vulns = []htmlVuln{}
	for _, layer := range image.ImageLayers {
		for _, pkg := range layer.Packages {
			for _, vul := range pkg.Vulnerabilities {
				space := regexp.MustCompile(`\s+`)
				layerCreatedBy := space.ReplaceAllString(layer.CreatedBy, " ")

				newHtmlVuln := htmlVuln{
					CVE:               vul.Name,
					Severity:          strings.Title(vul.Severity),
					SeverityHTMLClass: vul.Severity,
					PkgName:           pkg.Name,
					PkgVersion:        pkg.Version,
					PkgFixed:          vul.FixVersion,
					Layer:             layerCreatedBy,
				}

				if score := vul.CVSSv3Score(); score != 0 {
					// CVSSv3
					newHtmlVuln.V3Score = score
					newHtmlVuln.UseV3Score = true
				} else if score = vul.CVSSv2Score(); score != 0 {
					// CVSSv2
					newHtmlVuln.V2Score = score
					newHtmlVuln.UseV2Score = true
				} else {
					// N/A
					newHtmlVuln.UseNoScore = true
				}

				vulns = append(vulns, newHtmlVuln)
			}
		}
	}

	// order by severity
	sort.Slice(vulns, func(i, j int) bool {
		return severityOrder(vulns[i].Severity) < severityOrder(vulns[j].Severity)
	})

	// add the row height after ordering
	for row := range vulns {
		vulns[row].RowHeight = (int32(row) * magicTableHeightMultiplier) + magicTableHeightMultiplier
	}

	return vulns
}
