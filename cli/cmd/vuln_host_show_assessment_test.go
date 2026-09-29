package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lacework/go-sdk/v2/api"
	"github.com/lacework/go-sdk/v2/internal/capturer"
	"github.com/lacework/go-sdk/v2/internal/lacework"
)

func TestGetUniqueHostEvalGUID(t *testing.T) {
	expectedEvalGUID := "12345"
	actualEvalGUID := getUniqueHostEvalGUID(mockVulnerabilitiesHostResponse(expectedEvalGUID))

	assert.Equal(t, expectedEvalGUID, actualEvalGUID)
}

func mockVulnerabilitiesHostResponse(evalGUID string) api.VulnerabilitiesHostResponse {
	return api.VulnerabilitiesHostResponse{
		Data: []api.VulnerabilityHost{
			{
				EvalGUID:  "54321",
				StartTime: time.Now().AddDate(0, 0, -2),
			},
			{
				EvalGUID:  "54321",
				StartTime: time.Now().AddDate(0, 0, -2),
			},
			{
				EvalGUID:  evalGUID,
				StartTime: time.Now(),
			},
			{
				EvalGUID:  "98765",
				StartTime: time.Now().AddDate(0, 0, -1),
			},
			{
				EvalGUID:  evalGUID,
				StartTime: time.Now(),
			},
		},
	}
}

const agentOnlyMid = "42"

// newAgentOnlyHostAPI serves a host that has only Agent assessments, answering each search by the
// evalCtx.collector_type it filters on, and counts the searches per collector type.
func newAgentOnlyHostAPI(t *testing.T) map[string]int {
	t.Helper()
	var (
		mu    sync.Mutex
		calls = map[string]int{}
	)
	server := lacework.MockServer()
	server.MockToken("TOKEN")
	server.MockAPI("Entities/Machines/search", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, `{"data": [{"mid": %s, "hostname": "agent-only"}], "paging": {}}`, agentOnlyMid)
	})
	server.MockAPI("Vulnerabilities/Hosts/search", func(w http.ResponseWriter, r *http.Request) {
		var filter api.SearchFilter
		require.NoError(t, json.NewDecoder(r.Body).Decode(&filter))
		collector := ""
		for _, f := range filter.Filters {
			if f.Field == "evalCtx.collector_type" {
				collector = f.Value
			}
		}
		mu.Lock()
		calls[collector]++
		mu.Unlock()

		if collector != vulnHostCollectorTypeAgent {
			fmt.Fprint(w, `{"data": [], "paging": {}}`)
			return
		}
		fmt.Fprintf(w, `{"data": [{
			"mid": %s, "evalGuid": "eval-agent", "startTime": %q, "vulnId": "CVE-2026-0001",
			"severity": "High", "status": "Active",
			"evalCtx": {"collector_type": "Agent", "hostname": "agent-only"},
			"featureKey": {"name": "openssl", "version_installed": "1.0"}
		}], "paging": {}}`, agentOnlyMid, time.Now().UTC().Format(time.RFC3339))
	})
	t.Cleanup(server.Close)

	client, err := api.NewClient("test", api.WithToken("TOKEN"), api.WithURL(server.URL()))
	require.NoError(t, err)
	prevAPI, prevCache, prevNoCache, prevState := cli.LwApi, cli.Cache, cli.noCache, vulCmdState
	cli.LwApi = client
	cli.InitCache(t.TempDir())
	cli.noCache = false
	cli.NonInteractive()
	cli.EnableJSONOutput()
	t.Cleanup(func() {
		cli.LwApi, cli.Cache, cli.noCache, vulCmdState = prevAPI, prevCache, prevNoCache, prevState
		cli.Interactive()
		cli.EnableHumanOutput()
		vulHostShowAssessmentCmd.Flags().Lookup("collector_type").Changed = false
	})
	return calls
}

// showAssessment runs the command as `show-assessment <mid> [--collector_type <collector>]`; an
// empty collector leaves the flag unset.
func showAssessment(collector string) (string, error) {
	flag := vulHostShowAssessmentCmd.Flags().Lookup("collector_type")
	flag.Changed = collector != ""
	vulCmdState.CollectorType = collector
	if collector == "" {
		vulCmdState.CollectorType = vulnHostCollectorTypeAgentless
	}

	var err error
	out := capturer.CaptureOutput(func() {
		err = vulHostShowAssessmentCmd.RunE(vulHostShowAssessmentCmd, []string{agentOnlyMid})
	})
	return out, err
}

// The cached Agent assessment used to answer the Agentless request too, since the cache key was the
// machine alone (LINK-4264).
func TestShowAssessmentCachesPerCollectorType(t *testing.T) {
	calls := newAgentOnlyHostAPI(t)

	out, err := showAssessment(vulnHostCollectorTypeAgent)
	require.NoError(t, err)
	assert.Contains(t, out, "CVE-2026-0001")

	_, err = showAssessment(vulnHostCollectorTypeAgentless)
	assert.ErrorContains(t, err, "no data found with Agentless collector")
	assert.Equal(t, 1, calls[vulnHostCollectorTypeAgentless], "the Agentless request reaches the API")
}

func TestShowAssessmentServesARepeatFromTheCache(t *testing.T) {
	calls := newAgentOnlyHostAPI(t)

	_, err := showAssessment(vulnHostCollectorTypeAgent)
	require.NoError(t, err)
	agentCalls := calls[vulnHostCollectorTypeAgent]

	out, err := showAssessment(vulnHostCollectorTypeAgent)
	require.NoError(t, err)
	assert.Contains(t, out, "CVE-2026-0001")
	assert.Equal(t, agentCalls, calls[vulnHostCollectorTypeAgent], "the repeat is served from the cache")
}

// Asking for Agentless explicitly used to fall back to Agent too, returning Agent data for it.
func TestShowAssessmentFallsBackToAgentOnlyByDefault(t *testing.T) {
	newAgentOnlyHostAPI(t)

	out, err := showAssessment("")
	require.NoError(t, err, "with no --collector_type, a host without Agentless data falls back to Agent")
	assert.Contains(t, out, "CVE-2026-0001")

	_, err = showAssessment(vulnHostCollectorTypeAgentless)
	assert.ErrorContains(t, err, "no data found with Agentless collector",
		"an explicit Agentless request gets neither the fallback nor its cached result")
}

// With nothing cached, an explicit Agentless request used to fall back to Agent and return the Agent
// assessment as if it were Agentless.
func TestShowAssessmentHonoursAnExplicitAgentless(t *testing.T) {
	calls := newAgentOnlyHostAPI(t)

	_, err := showAssessment(vulnHostCollectorTypeAgentless)

	assert.ErrorContains(t, err, "no data found with Agentless collector")
	assert.Zero(t, calls[vulnHostCollectorTypeAgent], "an explicit Agentless request never asks for Agent data")
}
