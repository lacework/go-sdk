//go:build !windows && generation

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/lacework/go-sdk/lwgenerate/aws_eks_audit"

	"github.com/Netflix/go-expect"
	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/stretchr/testify/assert"
)

func assertEksAuditTerraformSaved(t *testing.T, message string) {
	assert.Contains(t, message, "Terraform code saved in")
}

const (
	eksPath = "/lacework/aws_eks_audit/"
)

func runEksAuditGenerateTest(t *testing.T, conditions func(*expect.Console), args ...string) string {
	os.Setenv("HOME", tfPath)

	hcl_path := filepath.Join(tfPath, eksPath, "main.tf")

	runFakeTerminalTestFromDir(t, tfPath, conditions, args...)
	out, err := os.ReadFile(hcl_path)
	if err != nil {
		return fmt.Sprintf("main.tf not found: %s", err)
	}

	t.Cleanup(func() {
		os.Remove(hcl_path)
	})

	result := terraformValidate(filepath.Join(tfPath, eksPath))

	assert.True(t, result.Valid)

	return string(out)
}

func TestGenerationEksSingleRegion(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
	)

	assertEksAuditTerraformSaved(t, final)

	regionClusterMap := make(map[string][]string)
	regionClusterMap["us-west-1"] = []string{"cluster1", "cluster2"}
	buildTf, _ := aws_eks_audit.NewTerraform(aws_eks_audit.WithParsedRegionClusterMap(regionClusterMap)).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationEksSingleRegionAdvancedBucket(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "y"},
				MsgMenu{cmd.EksAuditConfigureBucket, 0},
				MsgRsp{cmd.QuestionUseExistingBucket, "n"},
				MsgRsp{cmd.QuestionEksAuditBucketVersioning, "y"},
				MsgRsp{cmd.QuestionEksAuditMfaDeleteS3Bucket, "y"},
				MsgRsp{cmd.QuestionEksAuditBucketEncryption, "y"},
				MsgRsp{cmd.QuestionEksAuditBucketExistingKey, "n"},
				MsgRsp{cmd.QuestionEksAuditBucketSseAlgorithm, ""},
				MsgRsp{cmd.QuestionEksAuditKmsKeyRotation, "y"},
				MsgRsp{cmd.QuestionEksAuditKmsKeyDeletionDays, "30"},
				MsgRsp{cmd.QuestionEksAuditBucketLifecycle, "30"},
				MsgRsp{cmd.QuestionEksAuditAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
	)

	assertEksAuditTerraformSaved(t, final)

	regionClusterMap := make(map[string][]string)
	regionClusterMap["us-west-1"] = []string{"cluster1", "cluster2"}
	buildTf, _ := aws_eks_audit.NewTerraform(
		aws_eks_audit.WithParsedRegionClusterMap(regionClusterMap),
		aws_eks_audit.EnableBucketVersioning(true),
		aws_eks_audit.EnableBucketMfaDelete(),
		aws_eks_audit.EnableBucketEncryption(true),
		aws_eks_audit.EnableKmsKeyRotation(true),
		aws_eks_audit.WithKmsKeyDeletionDays(30),
		aws_eks_audit.WithBucketLifecycleExpirationDays(30),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationEksSingleRegionAdvancedBucketExistingKey(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "y"},
				MsgMenu{cmd.EksAuditConfigureBucket, 0},
				MsgRsp{cmd.QuestionUseExistingBucket, "n"},
				MsgRsp{cmd.QuestionEksAuditBucketVersioning, "y"},
				MsgRsp{cmd.QuestionEksAuditMfaDeleteS3Bucket, "y"},
				MsgRsp{cmd.QuestionEksAuditBucketEncryption, "y"},
				MsgRsp{cmd.QuestionEksAuditBucketExistingKey, "y"},
				MsgRsp{cmd.QuestionEksAuditBucketSseAlgorithm, ""},
				MsgRsp{cmd.QuestionEksAuditBucketKeyArn, "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"},
				MsgRsp{cmd.QuestionEksAuditBucketLifecycle, "30"},
				MsgRsp{cmd.QuestionEksAuditAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
	)

	assertEksAuditTerraformSaved(t, final)

	regionClusterMap := make(map[string][]string)
	regionClusterMap["us-west-1"] = []string{"cluster1", "cluster2"}
	buildTf, _ := aws_eks_audit.NewTerraform(
		aws_eks_audit.WithParsedRegionClusterMap(regionClusterMap),
		aws_eks_audit.EnableBucketVersioning(true),
		aws_eks_audit.EnableBucketMfaDelete(),
		aws_eks_audit.EnableBucketEncryption(true),
		aws_eks_audit.WithBucketSseKeyArn("arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"),
		aws_eks_audit.WithBucketLifecycleExpirationDays(30),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationEksSingleRegionAdvancedCrossAccountIam(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "y"},
				MsgMenu{cmd.EksAuditExistingCaIamRole, 1},
				MsgRsp{cmd.QuestionEksAuditExistingCaIamArn, "arn:aws:iam::249446771485:role/2537e820-ca-role"},
				MsgRsp{cmd.QuestionEksAuditExistingCaIamExtID, "123456789"},
				MsgRsp{cmd.QuestionEksAuditAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
	)

	assertEksAuditTerraformSaved(t, final)

	regionClusterMap := make(map[string][]string)
	regionClusterMap["us-west-1"] = []string{"cluster1", "cluster2"}
	buildTf, _ := aws_eks_audit.NewTerraform(
		aws_eks_audit.WithParsedRegionClusterMap(regionClusterMap),
		aws_eks_audit.WithExistingCrossAccountIamRole(
			&aws_eks_audit.ExistingCrossAccountIamRoleDetails{
				Arn:        "arn:aws:iam::249446771485:role/2537e820-ca-role",
				ExternalId: "123456789",
			}),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationEksSingleRegionAdvancedFirehoseSettings(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "y"},
				MsgMenu{cmd.EksAuditConfigureFh, 2},
				MsgRsp{cmd.QuestionEksAuditExistingFhIamRole, "y"},
				MsgRsp{cmd.QuestionEksAuditExistingFhIamArn, "arn:aws:iam::249446771485:role/2537e820-fh-role"},
				MsgRsp{cmd.QuestionEksAuditFhEncryption, "y"},
				MsgRsp{cmd.QuestionEksAuditFhEncryptionKeyArn, ""},
				MsgRsp{cmd.QuestionEksAuditAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
	)

	assertEksAuditTerraformSaved(t, final)

	regionClusterMap := make(map[string][]string)
	regionClusterMap["us-west-1"] = []string{"cluster1", "cluster2"}
	buildTf, _ := aws_eks_audit.NewTerraform(
		aws_eks_audit.WithParsedRegionClusterMap(regionClusterMap),
		aws_eks_audit.WithExistingFirehoseIamRoleArn(
			"arn:aws:iam::249446771485:role/2537e820-fh-role",
		),
		aws_eks_audit.EnableFirehoseEncryption(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationEksSingleRegionAdvancedExistingCloudwatchRole(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "y"},
				MsgMenu{cmd.EksAuditExistingCwIamRole, 3},
				MsgRsp{cmd.QuestionEksAuditExistingCwIamArn, "arn:aws:iam::249446771485:role/2537e820-cw-role"},
				MsgRsp{cmd.QuestionEksAuditAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
	)

	assertEksAuditTerraformSaved(t, final)

	regionClusterMap := make(map[string][]string)
	regionClusterMap["us-west-1"] = []string{"cluster1", "cluster2"}
	buildTf, _ := aws_eks_audit.NewTerraform(
		aws_eks_audit.WithParsedRegionClusterMap(regionClusterMap),
		aws_eks_audit.WithExistingCloudWatchIamRoleArn(
			"arn:aws:iam::249446771485:role/2537e820-cw-role",
		),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationEksSingleRegionAdvancedSnsSettings(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "y"},
				MsgMenu{cmd.EksAuditConfigureSns, 4},
				MsgRsp{cmd.QuestionEksAuditSnsEncryption, "y"},
				MsgRsp{cmd.QuestionEksAuditSnsEncryptionKeyArn, "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"},
				MsgRsp{cmd.QuestionEksAuditAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
	)

	assertEksAuditTerraformSaved(t, final)

	regionClusterMap := make(map[string][]string)
	regionClusterMap["us-west-1"] = []string{"cluster1", "cluster2"}
	buildTf, _ := aws_eks_audit.NewTerraform(
		aws_eks_audit.WithParsedRegionClusterMap(regionClusterMap),
		aws_eks_audit.EnableSnsTopicEncryption(true),
		aws_eks_audit.WithSnsTopicEncryptionKeyArn(
			"arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903",
		),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationEksSingleRegionAdvancedCustomIntegrationName(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "y"},
				MsgMenu{cmd.EksAuditIntegrationNameOpt, 5},
				MsgRsp{cmd.QuestionEksAuditCustomIntegrationName,
					"custom eks audit integration name"},
				MsgRsp{cmd.QuestionEksAuditAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
	)

	assertEksAuditTerraformSaved(t, final)

	regionClusterMap := make(map[string][]string)
	regionClusterMap["us-west-1"] = []string{"cluster1", "cluster2"}
	buildTf, _ := aws_eks_audit.NewTerraform(
		aws_eks_audit.WithParsedRegionClusterMap(regionClusterMap),
		aws_eks_audit.WithEksAuditIntegrationName("custom eks audit integration name"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationEksSingleRegionAdvancedCustomOutputLocation(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	dir, err := os.MkdirTemp("", "t")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "y"},
				MsgMenu{cmd.EksAuditAdvancedOptLocation, 6},
				MsgRsp{cmd.QuestionEksAuditCustomizeOutputLocation, dir},
				MsgRsp{cmd.QuestionEksAuditAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
	)

	assertEksAuditTerraformSaved(t, final)

	result, _ := os.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/main.tf", dir)))

	regionClusterMap := make(map[string][]string)
	regionClusterMap["us-west-1"] = []string{"cluster1", "cluster2"}
	buildTf, _ := aws_eks_audit.NewTerraform(
		aws_eks_audit.WithParsedRegionClusterMap(regionClusterMap),
	).Generate()
	assert.Equal(t, buildTf, string(result))
}

func TestGenerationEksSingleRegionExistingTerraform(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

	dir, err := os.MkdirTemp("", "t")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	if err := os.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/main.tf", dir)), []byte{}, 0644); err != nil {
		panic(err)
	}

	runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "y"},
				MsgMenu{cmd.EksAuditAdvancedOptLocation, 6},
				MsgRsp{cmd.QuestionEksAuditCustomizeOutputLocation, dir},
				MsgRsp{cmd.QuestionEksAuditAnotherAdvancedOpt, "n"},
				MsgRsp{fmt.Sprintf("%s/main.tf already exists, overwrite?", dir), "n"},
			})
		},
		"generate",
		"k8s",
		"eks",
	)

	// Ensure CLI ran correctly
	data, err := os.ReadFile(fmt.Sprintf("%s/main.tf", dir))
	if err != nil {
		panic(err)
	}

	assert.Empty(t, data)
}

func TestGenerateEksPrefix(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	prefix := "prefix-"

	tfResult := runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
		"--prefix",
		prefix,
	)

	assertEksAuditTerraformSaved(t, final)

	regionClusterMap := make(map[string][]string)
	regionClusterMap["us-west-1"] = []string{"cluster1", "cluster2"}
	buildTf, _ := aws_eks_audit.NewTerraform(
		aws_eks_audit.WithParsedRegionClusterMap(regionClusterMap),
		aws_eks_audit.WithPrefix(prefix),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationEksOverwriteOutput(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	dir := createDummyTOMLConfig()
	defer os.RemoveAll(dir)

	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", homeCache)

	output_dir := createDummyTOMLConfig()
	defer os.RemoveAll(output_dir)

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
		"--output",
		output_dir,
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s", output_dir))

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "n"},
				MsgRsp{"already exists, overwrite?", "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
		"--output",
		output_dir,
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s", output_dir))
}

func TestGenerationEksMultiRegion(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "y"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditAdditionalRegion, "y"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-east-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster3"},
				MsgRsp{cmd.QuestionEksAuditAdditionalRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
	)

	assertEksAuditTerraformSaved(t, final)

	regionClusterMap := make(map[string][]string)
	regionClusterMap["us-west-1"] = []string{"cluster1", "cluster2"}
	regionClusterMap["us-east-1"] = []string{"cluster3"}
	buildTf, _ := aws_eks_audit.NewTerraform(aws_eks_audit.WithParsedRegionClusterMap(regionClusterMap)).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationEksMultiRegionCliFlag(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	regionClusterMap := make(map[string][]string)
	regionClusterMap["us-west-1"] = []string{"cluster1", "cluster2"}
	regionClusterMap["us-east-1"] = []string{"cluster3"}

	tfResult := runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{fmt.Sprintf(
					cmd.QuestionEksAuditRegionClusterCurrent,
					regionClusterMap,
				), "n"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
		"--region_clusters",
		"us-west-1=cluster1,cluster2",
		"--region_clusters",
		"us-east-1=cluster3",
	)

	assertEksAuditTerraformSaved(t, final)

	buildTf, _ := aws_eks_audit.NewTerraform(aws_eks_audit.WithParsedRegionClusterMap(regionClusterMap)).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationEksNonInteractive(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
		"--region_clusters",
		"us-west-1=cluster1,cluster2",
		"--noninteractive",
	)

	assertEksAuditTerraformSaved(t, final)

	regionClusterMap := make(map[string][]string)
	regionClusterMap["us-west-1"] = []string{"cluster1", "cluster2"}

	buildTf, _ := aws_eks_audit.NewTerraform(aws_eks_audit.WithParsedRegionClusterMap(regionClusterMap)).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationSuppliedBucketArn(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runEksAuditGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEksAuditMultiRegion, "n"},
				MsgRsp{cmd.QuestionEksAuditRegion, "us-west-1"},
				MsgRsp{cmd.QuestionEksAuditRegionClusters, "cluster1,cluster2"},
				MsgRsp{cmd.QuestionEksAuditConfigureAdvanced, "y"},
				MsgMenu{cmd.EksAuditConfigureBucket, 0},
				MsgRsp{cmd.QuestionUseExistingBucket, "y"},
				MsgRsp{cmd.QuestionExistingBucketArn, "arn:aws:s3:::bucket-name"},
				MsgRsp{cmd.QuestionEksAuditAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"eks",
	)

	assertEksAuditTerraformSaved(t, final)

	regionClusterMap := make(map[string][]string)
	regionClusterMap["us-west-1"] = []string{"cluster1", "cluster2"}
	buildTf, _ := aws_eks_audit.NewTerraform(
		aws_eks_audit.WithParsedRegionClusterMap(regionClusterMap),
		aws_eks_audit.EnableUseExistingBucket(),
		aws_eks_audit.WithExistingBucketArn("arn:aws:s3:::bucket-name"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}
