package main

import (
	"fmt"

	"github.com/lacework/go-sdk/lwpager"
)

func main() {
	pager, err := lwpager.Start()
	if err != nil {
		panic(err)
	}

	defer pager.Wait()

	// Output inside less:
	// Hi, this test is being paged!
	// (END)

	fmt.Fprintln(pager.Out, "Hi, this test is being paged!")
}
