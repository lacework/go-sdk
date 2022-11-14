package resources

import "time"

type InstanceDetails struct {
	InstanceID string    `json:"INSTANCE_ID"`
	Type       string    `json:"TYPE"`
	State      string    `json:"STATE"`
	Name       string    `json:"NAME"`
	Region     string    `json:"REGION"`
	Zone       string    `json:"ZONE"`
	ImageID    string    `json:"IMAGE_ID"`
	VpcID      string    `json:"VPC_ID"`
	AccountID  string    `json:"ACCOUNT_ID"`
	PublicIP   string    `json:"PUBLIC_IP"`
	PrivateIP  string    `json:"PRIVATE_IP"`
	LaunchTime time.Time `json:"LAUNCH_TIME"`

	// HardwareDetails contains the "expanded out" CPU and Memory information for the instance type.
	HardwareDetails HardwareDetails `json:"HARDWARE_DETAILS"`

	Tags        map[string]string `json:"TAGS"`
	Props       map[string]string `json:"PROPS"`
	ConfigProps ConfigProps       `json:"CONFIG_PROPS"`
}

type HardwareDetails struct {
	CpuVendorID      string `json:"CPU_VENDOR_ID"`
	CpuModelName     string `json:"CPU_MODEL_NAME"`
	CpuCores         int    `json:"CPU_CORES"`
	MemoryTotalBytes int64  `json:"MEMORY_TOTAL_BYTES"`
	MemoryECCType    string `json:"MEMORY_ECC_TYPE"`
}

type ConfigProps struct {
	ScanFrequency           int64 `json:"SCAN_FREQUENCY"`
	ScanContainers          bool  `json:"SCAN_CONTAINERS"`
	ScanHostVulnerabilities bool  `json:"SCAN_HOST_VULNERABILITIES"`
	// Deprecated options:
	// ScanImages bool  `json:"SCAN_IMAGES"`

	// Cross-account role used for org access and internal access.
	CrossAccountCredentials ConfigCredentialsProps `json:"CROSS_ACCOUNT_CREDENTIALS"`

	// Org-specific properties
	ManagementAccount string `json:"MANAGEMENT_ACCOUNT"`
	MonitoredAccounts string `json:"MONITORED_ACCOUNTS"`
	ScanningAccount   string `json:"SCANNING_ACCOUNT"`
}

type ConfigCredentialsProps struct {
	RoleARN    string `json:"ROLE_ARN"`
	ExternalId string `json:"EXTERNAL_ID"`
}
