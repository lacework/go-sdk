package integration

func init() {
	lw, err := laceworkIntegrationTestClient()
	if err == nil {
		// Check if the default machineID(30) exists
		if !machineIDExists(lw) {
			setHostVulnTestMachineID(lw)
		}
		// set id for report definition tests
		setReportDefinition(lw)
	}
}
