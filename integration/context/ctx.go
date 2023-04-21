package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/lacework/go-sdk/internal/array"
)

var (
	ctx     IntegrationCtx
	allTags []string
)

func init() {
	ex, err := os.Getwd()
	if err != nil {
		log.Fatal("unable to determine ctx.toml path", err)
	}
	path := fmt.Sprintf("%s/integration/context/ctx.toml", ex)
	log.Print(path)
	if _, err := toml.DecodeFile(path, &ctx); err != nil {
		log.Fatal("unable to decode integration ctx config")
	}

	for tag := range ctx {
		allTags = append(allTags, tag)
	}

}
func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Println("no args provided")
		os.Exit(0)
	}

	testContext := integrationCtx(args)
	if len(testContext) == 0 {
		log.Println("no matching tags found for changes")
		os.Exit(0)
	}

	tags := array.Unique(testContext)

	// if 'all' is selected return all tags
	if array.ContainsStr(tags, "all") {
		log.Printf("running all integration tests: %s", strings.Join(allTags, ","))
		fmt.Print(strings.Join(allTags, " "))
		return
	}

	log.Printf("determined test context tags: %s", strings.Join(tags, ","))
	fmt.Print(strings.Join(tags, " "))
}

func integrationCtx(args []string) (buildTags []string) {
	log.Print("determining context...")
	for _, file := range args {
		for tagKey, tag := range ctx {
			// handle matching files
			if file != "" && array.ContainsStr(tag.Files, file) || array.ContainsStr(tag.Dirs, file) {
				log.Println(file)
				if !array.ContainsStr(buildTags, tagKey) {
					buildTags = append(buildTags, tagKey)
					log.Printf("added tag %s", tagKey)
				}
				continue
			}
		}
	}
	return
}

type IntegrationCtx map[string]TestCtx

type TestCtx struct {
	Files []string `toml:"files"`
	Dirs  []string `toml:"dirs"`
}
