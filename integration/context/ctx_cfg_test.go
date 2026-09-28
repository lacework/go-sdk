package main

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationCtx(t *testing.T) {
	t.Chdir("../..")

	for _, tc := range []struct {
		name  string
		files []string
		want  []string
	}{
		{"a command maps to its tag", []string{"cli/cmd/preflight_aws.go"}, []string{"help", "preflight"}},
		{"a directory maps to its tag", []string{"lwpreflight/aws/detail.go"}, []string{"help", "preflight"}},
		{"a changed test runs its own tags", []string{"integration/aws_generation_test.go"}, []string{"generation", "help"}},
		{"docs change nothing", []string{"README.md", "docs/README.md"}, nil},
		{"the argument separator is ignored", []string{"--", "README.md"}, nil},
		{"an unmapped file runs every tag", []string{"cli/cmd/preflight_aws.go", "api/client.go"}, allTags()},
		{"an [all] entry runs every tag", []string{"cli/cmd/root.go"}, allTags()},
		{"a same-named file outside cli/cmd is not the command", []string{"lwgenerate/aws/account.go"}, []string{"generation", "help"}},
		{"an untagged test file runs every tag", []string{"integration/framework_test.go"}, allTags()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, integrationCtx(tc.files))
		})
	}
}

// The Makefile shards INTEGRATION_TEST_TAGS in CI and ctx.toml selects among them, so a tag in only
// one of the two either never runs on a PR or is selected and never run.
func TestCtxTagsMatchTheMakefile(t *testing.T) {
	makefile, err := os.ReadFile("../../Makefile")
	require.NoError(t, err)
	block := regexp.MustCompile(`(?s)INTEGRATION_TEST_TAGS=(.*?)\n\n`).FindSubmatch(makefile)
	require.NotNil(t, block, "INTEGRATION_TEST_TAGS not found in the Makefile")
	makeTags := strings.Fields(strings.ReplaceAll(string(block[1]), `\`, " "))
	slices.Sort(makeTags)

	assert.Equal(t, makeTags, allTags())
}

// A stale entry silently stops mapping its changes to the tag.
func TestCtxEntriesExist(t *testing.T) {
	t.Chdir("../..")

	for tag, tagCtx := range ctx {
		for _, file := range tagCtx.Files {
			assert.FileExists(t, filepath.Join("cli/cmd", file), "ctx.toml [%s]", tag)
		}
		for _, dir := range tagCtx.Dirs {
			assert.DirExists(t, dir, "ctx.toml [%s]", tag)
		}
	}
}

// An untagged test file compiles into every shard, so its tests run once per shard in parallel
// against the same account.
func TestIntegrationTestsAreTagged(t *testing.T) {
	files, err := filepath.Glob("../*_test.go")
	require.NoError(t, err)
	require.NotEmpty(t, files)

	for _, file := range files {
		src, err := os.ReadFile(file)
		require.NoError(t, err)
		// TestMain takes *testing.M and is shared setup every shard needs, so it stays untagged.
		if !regexp.MustCompile(`(?m)^func Test\w*\(t \*testing\.T\)`).Match(src) {
			continue
		}
		assert.Regexp(t, `(?m)^//go:build `, string(src), "%s has tests but no build tag", filepath.Base(file))
	}
}
