package version

var (
	// All the following "unknown" variables are being injected at
	// build time via the cross-platform directive inside the Makefile
	//
	// Version is the semver coming from the VERSION file
	Version = "unknown"

	// GitSHA is the git ref that the cli was built from
	GitSHA = "unknown"

	// BuildTime is a human-readable time when the cli was built at
	BuildTime = "unknown"
)
