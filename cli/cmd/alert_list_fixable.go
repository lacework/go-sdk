package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
)

const remediateComponentName string = "remediate"

// isRemediateInstalled returns true if the remediate component is installed
func (c *cliState) isRemediateInstalled() bool {
	return c.IsComponentInstalled(remediateComponentName)
}

// getTemplateIdentifiers runs the remediate component to retrieve a list
// of remediation template identifiers
func getRemediationTemplateIDs() ([]string, error) {
	remediate, found := cli.LwComponents.GetComponent(remediateComponentName)
	if !found {
		return []string{}, errors.New("remediate component not found")
	}

	// set up environment variables
	envs := []string{
		fmt.Sprintf("LW_COMPONENT_NAME=%s", remediateComponentName),
		"LW_JSON=true",
		"LW_NONINTERACTIVE=true",
	}
	for _, e := range cli.envs() {
		// don't let LW_JSON / LW_NONINTERACTIVE through here
		if strings.HasPrefix(e, "LW_JSON=") || strings.HasPrefix(e, "LW_NONINTERACTIVE=") {
			continue
		}
		envs = append(envs, e)
	}
	stdout, stderr, err := remediate.RunAndReturn([]string{"ls", "templates"}, nil, envs...)
	if err != nil {
		cli.Log.Debugw("remediate error details", "stderr", stderr)
		return []string{}, err
	}

	var templates []map[string]interface{}
	err = json.Unmarshal([]byte(stdout), &templates)
	if err != nil {
		return []string{}, err
	}

	templateIDs := []string{}
	for _, template := range templates {
		v, ok := template["id"]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		templateIDs = append(templateIDs, s)
	}
	return templateIDs, nil
}

// filterFixableAlerts identifies which alerts have corresponding remediation template IDs
// and returns those which don't
func filterFixableAlerts(
	alerts api.Alerts,
	getIDsFunc func() ([]string, error),
) (api.Alerts, error) {
	templateIDs, err := getIDsFunc()
	if err != nil {
		return alerts, err
	}

	fixableAlerts := api.Alerts{}
	for _, alert := range alerts {
		if alert.PolicyID == "" {
			continue
		}
		found := false
		// Historically alerts did not consistently populate policyID and
		// templates were named arbitrarily.
		// If and when policies explicitly reference templates we will no longer need
		// any inference logic.
		for _, id := range templateIDs {
			if id == alert.PolicyID {
				fixableAlerts = append(fixableAlerts, alert)
				found = true
				break
			}
		}
		if found {
			continue
		}
		// Another interesting problem that we have is that policyIDs are dynamic
		// For instance, on dev7 policy lwcustom-11 is dev7-lwcustom-11
		// On some other environment it might be someother-lwcustom-11
		dynamicIDRE := regexp.MustCompile(`^\w+-\d+$`)
		// Iterate through the templates looking for those with dynamic policy IDs
		for _, id := range templateIDs {
			if dynamicIDRE.MatchString(id) {
				// if the policyID of the alert ends with -<id>
				// i.e. if dev7-lwcustom-11 endswith -lwcustom-11
				if strings.HasSuffix(alert.PolicyID, fmt.Sprintf("-%s", id)) {
					fixableAlerts = append(fixableAlerts, alert)
					break
				}
			}
		}
	}
	return fixableAlerts, nil
}
