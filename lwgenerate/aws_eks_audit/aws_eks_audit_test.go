package aws_eks_audit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper for combining string expected values
//func reqProviderAndRegion(extraInputs ...string) string {
//	base := requiredProviders + "\n" + awsProvider
//	countInputs := len(extraInputs)
//	for i, e := range extraInputs {
//		if i < countInputs {
//			base = base + "\n" + e
//		}
//
//		if i >= countInputs {
//			base = base + e
//		}
//	}
//	return base
//}

func TestGenerationCloudTrail(t *testing.T) {
	clusterMap := make(map[string][]string)
	regionOne := []string{"cluster1", "cluster2"}
	regionTwo := []string{"cluster3"}
	clusterMap["us-east-1"] = regionOne
	clusterMap["us-east-2"] = regionTwo
	hcl, err := NewTerraform(WithRegionClusterMap(clusterMap)).Generate()
	print(hcl)
	assert.Nil(t, err)

}
