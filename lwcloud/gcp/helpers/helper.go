package helpers

import (
	"fmt"
	"os"
	"strings"
)

const (
	GCP_PROJECT_TYPE      = "PROJECT"
	GCP_ORGANIZATION_TYPE = "ORGANIZATION"
)

func SkipEntry(Entry string, skipList, allowList map[string]bool) bool {
	if skipList != nil {
		if _, skip := skipList[Entry]; skip {
			return true
		}
	}

	if allowList != nil {
		if _, allow := allowList[Entry]; allow {
			return false
		} else {
			// skip all other entries
			return true
		}
	}

	return false
}

func GetGcpFormatedLabel(in string) string {
	lower := strings.ToLower(in)
	out := strings.ReplaceAll(lower, ":", "-")
	out = strings.ReplaceAll(out, ".", "-")
	out = strings.ReplaceAll(out, "\"", "-")
	out = strings.ReplaceAll(out, "{", "-")
	out = strings.ReplaceAll(out, "}", "-")
	return out
}

func CombineErrors(old error, new error) error {
	if new == nil {
		return old
	}

	if old == nil {
		return new
	}

	return fmt.Errorf("%s, %s", old.Error(), new.Error())
}

func IsProjectScanScope() bool {
	return os.Getenv("GCP_SCAN_SCOPE") == GCP_PROJECT_TYPE
}

func IsOrgScanScope() bool {
	return os.Getenv("GCP_SCAN_SCOPE") == GCP_ORGANIZATION_TYPE
}
