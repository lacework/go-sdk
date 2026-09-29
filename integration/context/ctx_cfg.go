// Command ctx_cfg prints the integration test build tags that the given changed files need.
//
//	go run integration/context/ctx_cfg.go $(git diff --name-only origin/main...HEAD)
//
// Run it from the repository root: it reads integration test files to learn their build tags.
package main

import (
	_ "embed"
	"fmt"
	"go/build/constraint"
	"log"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
)

//go:embed ctx.toml
var ctxTOML string

var ctx IntegrationCtx

// skipped are changes no integration test can observe.
var skipped = struct{ suffixes, prefixes, files []string }{
	suffixes: []string{".md"},
	prefixes: []string{"docs/", ".github/ISSUE_TEMPLATE/", "integration/context/"},
	files:    []string{"CODEOWNERS", ".github/CODEOWNERS", "LICENSE", "VERSION"},
}

func init() {
	if _, err := toml.Decode(ctxTOML, &ctx); err != nil {
		log.Fatal("unable to decode integration ctx config: ", err)
	}
}

func main() {
	tags := integrationCtx(os.Args[1:])
	log.Printf("determined test context tags: %s", strings.Join(tags, ","))
	fmt.Print(strings.Join(tags, " "))
}

// integrationCtx maps changed files to build tags. A file no mapping claims runs every tag: a
// change nobody mapped is exactly the one whose blast radius is unknown.
func integrationCtx(files []string) []string {
	var tags []string
	for _, file := range files {
		// A unit test outside integration/ is never compiled into the CLI binary the integration
		// tests run, and the unit-test job already covers it.
		unitTest := strings.HasSuffix(file, "_test.go") && !strings.HasPrefix(file, "integration/")
		if file == "" || file == "--" || unitTest ||
			slices.ContainsFunc(skipped.suffixes, func(s string) bool { return strings.HasSuffix(file, s) }) ||
			slices.ContainsFunc(skipped.prefixes, func(p string) bool { return strings.HasPrefix(file, p) }) ||
			slices.Contains(skipped.files, file) {
			continue
		}

		var matched []string
		if strings.HasPrefix(file, "integration/") && strings.HasSuffix(file, "_test.go") {
			matched = testFileTags(file)
		} else {
			for tag, tagCtx := range ctx {
				// File entries name commands under cli/cmd/; directory entries cover everything below.
				if slices.Contains(tagCtx.Files, strings.TrimPrefix(file, "cli/cmd/")) && path.Dir(file) == "cli/cmd" ||
					slices.ContainsFunc(tagCtx.Dirs, func(dir string) bool { return strings.HasPrefix(file, dir) }) {
					matched = append(matched, tag)
				}
			}
		}
		if len(matched) == 0 || slices.Contains(matched, "all") {
			log.Printf("%s runs every integration test", file)
			return allTags()
		}
		tags = append(tags, matched...)
	}
	if len(tags) == 0 {
		return nil
	}

	// help tests are cheap and cover every command's wiring.
	tags = append(tags, "help")
	slices.Sort(tags)
	return slices.Compact(tags)
}

// testFileTags returns the integration tags named in a test file's //go:build line, so a changed
// test always runs itself. A deleted or untagged file yields none, which runs every tag.
func testFileTags(file string) []string {
	src, err := os.ReadFile(file)
	if err != nil {
		return nil
	}
	var tags []string
	// Walk the expression rather than Eval it: Eval short-circuits, so in "!windows && preflight"
	// it never reaches preflight.
	var walk func(constraint.Expr)
	walk = func(expr constraint.Expr) {
		switch e := expr.(type) {
		case *constraint.TagExpr:
			if _, ok := ctx[e.Tag]; ok && e.Tag != "all" {
				tags = append(tags, e.Tag)
			}
		case *constraint.NotExpr:
			walk(e.X)
		case *constraint.AndExpr:
			walk(e.X)
			walk(e.Y)
		case *constraint.OrExpr:
			walk(e.X)
			walk(e.Y)
		}
	}
	for _, line := range strings.Split(string(src), "\n") {
		if !constraint.IsGoBuild(line) {
			continue
		}
		expr, err := constraint.Parse(line)
		if err != nil {
			return nil
		}
		walk(expr)
	}
	return tags
}

func allTags() []string {
	var tags []string
	for tag := range ctx {
		if tag != "all" {
			tags = append(tags, tag)
		}
	}
	slices.Sort(tags)
	return tags
}

type IntegrationCtx map[string]TestCtx

type TestCtx struct {
	Files []string `toml:"files"`
	Dirs  []string `toml:"dirs"`
}
