//go:build policy && !windows

package integration

import (
	"os"
	"testing"

	"github.com/Netflix/go-expect"
)

func TestPolicyExceptionCreate(t *testing.T) {
	dir := createTOMLConfigFromCIvars()
	defer os.RemoveAll(dir)

	_ = runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectString(t, c, "Policy ID:")
			c.SendLine("\u001B[B")
			expectString(t, c, "Exception Description:")
			c.SendLine("TEST CLI")
			expectString(t, c, "Add constraint \"accountIds\"")
			c.SendLine("n")
			expectString(t, c, "Add constraint \"regionNames\"")
			c.SendLine("n")
			expectString(t, c, "Add constraint \"resourceNames\"")
			c.SendLine("n")
			expectString(t, c, "Add constraint \"resourceTags\"")
			c.SendLine("n")
			expectString(t, c, "policy exceptions must have at least 1 constraint")
		}, "policy-exception", "create",
	)
}
