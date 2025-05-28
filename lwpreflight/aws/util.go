package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
)

// ParseResourceName extracts the resource name from the caller identity Arn
// Examples:
//   - arn:aws:iam::123456789012:root -> root
//   - arn:aws:iam::123456789012:user/MyUser -> MyUser
//   - arn:aws:iam::123456789012:role/application_abc/component_xyz/RDSAccess -> RDSAccess
//   - arn:aws:iam::123456789012:assumed-role/preflight_ro/aws-go-sdk-00000000000 -> preflight_ro
func ParseResourceName(arnStr string) (string, error) {
	arnObj, err := arn.Parse(arnStr)
	if err != nil {
		return "", err
	}

	parts := strings.Split(arnObj.Resource, ":")
	lastStr := parts[len(parts)-1]
	paths := strings.Split(lastStr, "/")

	if strings.Contains(lastStr, "assumed-role") {
		return paths[len(paths)-2], nil
	}

	return paths[len(paths)-1], nil
}
