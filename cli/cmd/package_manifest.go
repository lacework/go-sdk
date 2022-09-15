//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/file"
)

var SupportedPackageManagers = []string{"dpkg-query", "rpm"} // @afiune can we support yum and apk?

type OS struct {
	Name    string
	Version string
}

var (
	procUAStatusFile = "/proc/1/root/var/lib/ubuntu-advantage/status.json"
	osReleaseFile    = "/etc/os-release"
	sysReleaseFile   = "/etc/system-release"
	rexNameFromID    = regexp.MustCompile(`^ID=(.*)$`)
	rexVersionID     = regexp.MustCompile(`^VERSION_ID=(.*)$`)
)

func (c *cliState) GeneratePackageManifest() (*api.VulnerabilitiesPackageManifest, error) {
	var (
		err   error
		start = time.Now()
	)

	defer func() {
		c.Event.DurationMs = time.Since(start).Milliseconds()
		if err == nil {
			// if this function returns an error, most likely,
			// the command will send a honeyvent with that error,
			// therefore we should duplicate events and only send
			// one here if there is NO error
			c.SendHoneyvent()
		}
	}()

	c.Event.Feature = featGenPkgManifest

	manifest := new(api.VulnerabilitiesPackageManifest)
	osInfo, err := c.GetOSInfo()
	if err != nil {
		return manifest, err
	}

	if osInfo.Name == "ubuntu" {
		// ESM support
		if c.IsEsmEnabled() {
			osInfo.Version += "esm"
		}
	}

	c.Event.AddFeatureField("os", osInfo.Name)
	c.Event.AddFeatureField("os_ver", osInfo.Version)

	manager, err := c.DetectPackageManager()
	if err != nil {
		return manifest, err
	}
	c.Event.AddFeatureField("pkg_manager", manager)

	var managerQuery []byte
	switch manager {
	case "rpm":
		managerQuery, err = exec.Command(
			"rpm", "-qa", "--queryformat", "%{NAME},%|EPOCH?{%{EPOCH}}:{0}|:%{VERSION}-%{RELEASE}\n",
		).Output()
		if err != nil {
			return manifest, errors.Wrap(err, "unable to query packages from package manager")
		}
	case "dpkg-query":
		managerQuery, err = exec.Command(
			"dpkg-query", "--show", "--showformat", "${Package},${Version}\n",
		).Output()
		if err != nil {
			return manifest, errors.Wrap(err, "unable to query packages from package manager")
		}
	case "yum":
		return manifest, errors.New("yum not yet supported")
	case "apk":
		apkInfo, err := exec.Command("apk", "info").Output()
		if err != nil {
			return manifest, errors.Wrap(err, "unable to query packages from package manager")
		}
		apkInfoT := strings.TrimSuffix(string(apkInfo), "\n")
		apkInfoArray := strings.Split(apkInfoT, "\n")

		apkInfoWithVersion, err := exec.Command("apk", "info", "-v").Output()
		if err != nil {
			return manifest, errors.Wrap(err, "unable to query packages from package manager")
		}
		apkInfoWithVersionT := strings.TrimSuffix(string(apkInfoWithVersion), "\n")
		apkInfoWithVersionArray := strings.Split(apkInfoWithVersionT, "\n")

		mq := []string{}
		for i, pkg := range apkInfoWithVersionArray {
			mq = append(mq,
				fmt.Sprintf("%s,%s",
					apkInfoArray[i],
					strings.Trim(
						strings.Replace(pkg, apkInfoArray[i], "", 1),
						"-",
					),
				),
			)
		}
		managerQuery = []byte(strings.Join(mq, "\n"))
	default:
		return manifest, errors.New(
			"this is most likely a mistake on us, please report it to support@lacework.com.",
		)
	}

	c.Log.Debugw("package-manager query", "raw", string(managerQuery))

	// @afiune this is an example of the output from the query we
	// send to the local package-manager:
	//
	// {PkgName},{PkgVersion}\n
	// ...
	// {PkgName},{PkgVersion}\n
	//
	// first, trim the last carriage return
	managerQueryOut := strings.TrimSuffix(string(managerQuery), "\n")
	// then, split by carriage return
	for _, pkg := range strings.Split(managerQueryOut, "\n") {
		// finally, split by comma to get PackageName and PackageVersion
		pkgDetail := strings.Split(pkg, ",")

		// the splitted package detail must be size of 2 elements
		if len(pkgDetail) != 2 {
			c.Log.Warnw("unable to parse package, expected length=2, skipping",
				"raw_pkg_details", pkg,
				"split_pkg_details", pkgDetail,
			)
			continue
		}

		manifest.OsPkgInfoList = append(manifest.OsPkgInfoList,
			api.VulnerabilitiesOsPkgInfo{
				Os:     osInfo.Name,
				OsVer:  osInfo.Version,
				Pkg:    pkgDetail[0],
				PkgVer: pkgDetail[1],
			},
		)
	}

	c.Event.AddFeatureField("total_manifest_pkgs", len(manifest.OsPkgInfoList))
	c.Log.Debugw("package-manifest", "raw", manifest)
	return c.removeInactivePackagesFromManifest(manifest, manager), nil
}

func (c *cliState) removeInactivePackagesFromManifest(manifest *api.VulnerabilitiesPackageManifest, manager string) *api.VulnerabilitiesPackageManifest {
	// Detect Active Kernel
	//
	// The default behavior of most linux distros is to keep the last N kernel packages
	// installed for users that need to fallback in case the new kernel do not boot.
	// However, the presence of the package does not mean that kernel is active.
	// We must continue to allow the standard kernel package preservation behavior
	// without providing false-positives of vulnerabilities that are not active.
	//
	// We will try to detect the active kernel and remove any other installed-inactive
	// kernel from the generated package manifest
	activeKernel, detected := c.detectActiveKernel()
	c.Event.AddFeatureField("active_kernel", activeKernel)
	if !detected {
		return manifest
	}

	newManifest := new(api.VulnerabilitiesPackageManifest)
	for i, pkg := range manifest.OsPkgInfoList {

		switch manager {
		case "rpm":
			kernelPkgName := "kernel"
			pkgVer := removeEpochFromPkgVersion(pkg.PkgVer)
			if pkg.Pkg == kernelPkgName && !strings.Contains(activeKernel, pkgVer) {
				// this package is NOT the active kernel
				c.Log.Warnw("inactive kernel package detected, removing from generated pkg manifest",
					"pkg_name", kernelPkgName,
					"pkg_version", pkg.PkgVer,
					"active_kernel", activeKernel,
				)
				c.Event.AddFeatureField(
					fmt.Sprintf("kernel_suppressed_%d", i),
					fmt.Sprintf("%s-%s", pkg.Pkg, pkg.PkgVer))
				continue
			}
		case "dpkg-query":
			kernelPkgName := "linux-image-"
			if strings.Contains(pkg.Pkg, kernelPkgName) {
				// this is a kernel package, trim the package name prefix to get the version
				kernelVer := strings.TrimPrefix(pkg.Pkg, kernelPkgName)

				if !strings.Contains(activeKernel, kernelVer) {
					// this package is NOT the active kernel
					c.Log.Warnw("inactive kernel package detected, removing from generated pkg manifest",
						"pkg_name", kernelPkgName,
						"pkg_version", pkg.PkgVer,
						"active_kernel", activeKernel,
					)
					c.Event.AddFeatureField(
						fmt.Sprintf("kernel_suppressed_%d", i),
						fmt.Sprintf("%s-%s", pkg.Pkg, pkg.PkgVer))
					continue
				}
			}
		}

		newManifest.OsPkgInfoList = append(newManifest.OsPkgInfoList, pkg)
	}

	if len(manifest.OsPkgInfoList) != len(newManifest.OsPkgInfoList) {
		c.Log.Debugw("package-manifest modified", "raw", newManifest)
	}
	return newManifest
}

func (c *cliState) detectActiveKernel() (string, bool) {
	kernel, err := exec.Command("uname", "-r").Output()
	if err != nil {
		c.Log.Warnw("unable to detect active kernel",
			"cmd", "uname -r",
			"error", err,
		)
		return "", false
	}
	return strings.TrimSuffix(string(kernel), "\n"), true
}

func (c *cliState) IsEsmEnabled() bool {
	type uaStatusFile struct {
		SchemaVersion string `json:"_schema_version,omitempty"`
		Status        string `json:"execution_status,omitempty"`
		Services      []struct {
			Name   string `json:"name,omitempty"`
			Status string `json:"status,omitempty"`
		} `json:"services,omitempty"`
	}

	if file.FileExists(procUAStatusFile) {
		c.Log.Debugw("detecting ubuntu ESM support", "file", procUAStatusFile)
		uaStatusBytes, err := ioutil.ReadFile(procUAStatusFile)
		if err != nil {
			c.Log.Warnw("unable to read UA status file", "error", err)
			return false
		}

		var uaStatus uaStatusFile
		if err = json.Unmarshal(uaStatusBytes, &uaStatus); err != nil {
			c.Log.Warnw("unable to unmarshal UA status file", "error", err)
			return false
		}

		for _, svc := range uaStatus.Services {
			if strings.Contains(svc.Name, "esm") && svc.Status == "enabled" {
				c.Log.Debug("ESM is enabled")
				return true
			}
		}

		c.Log.Debug("no UA service enabled")
		return false
	}

	c.Log.Warnw("unable to detect ubuntu ESM support, file not found", "file", procUAStatusFile)
	return false
}

func (c *cliState) GetOSInfo() (*OS, error) {
	osInfo := new(OS)

	c.Log.Debugw("detecting operating system information",
		"os", runtime.GOOS,
		"arch", runtime.GOARCH,
	)

	if file.FileExists(osReleaseFile) {
		c.Log.Debugw("parsing os release file", "file", osReleaseFile)
		return openOsReleaseFile(osReleaseFile)
	}

	if file.FileExists(sysReleaseFile) {
		c.Log.Debugw("parsing system release file", "file", sysReleaseFile)
		return openSystemReleaseFile(sysReleaseFile)
	}

	msg := `unsupported platform

For more information about supported platforms, visit:
  https://docs.lacework.com/host-vulnerability-assessment-overview`
	return osInfo, errors.New(msg)
}

func openSystemReleaseFile(filename string) (*OS, error) {
	osInfo := new(OS)

	f, err := os.Open(filename)

	if err != nil {
		return osInfo, err
	}

	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		m := strings.Split(s.Text(), " ")
		if len(m) > 0 {
			osInfo.Name = strings.ToLower(m[0])
			osInfo.Version = strings.ToLower(m[2])
			break
		}
	}

	return osInfo, err
}

func openOsReleaseFile(filename string) (*OS, error) {
	osInfo := new(OS)

	f, err := os.Open(filename)
	if err != nil {
		return osInfo, err
	}

	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		if m := rexNameFromID.FindStringSubmatch(s.Text()); m != nil {
			osInfo.Name = strings.Trim(m[1], `"`)
		} else if m := rexVersionID.FindStringSubmatch(s.Text()); m != nil {
			osInfo.Version = strings.Trim(m[1], `"`)
		}
	}

	return osInfo, err
}

func (c *cliState) DetectPackageManager() (string, error) {
	c.Log.Debugw("detecting package-manager")

	for _, manager := range SupportedPackageManagers {
		if c.checkPackageManager(manager) {
			c.Log.Debugw("detected", "package-manager", manager)
			return manager, nil
		}
	}
	msg := "unable to find supported package managers."
	msg = fmt.Sprintf("%s Supported package managers are %s.",
		msg, strings.Join(SupportedPackageManagers, ", "))
	return "", errors.New(msg)
}

func (c *cliState) checkPackageManager(manager string) bool {
	var (
		cmd    = exec.Command("which", manager)
		_, err = cmd.CombinedOutput()
	)
	if err != nil {
		c.Log.Debugw("error trying to check package-manager",
			"cmd", "which",
			"package-manager", manager,
			"error", err,
		)
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			return waitStatus.ExitStatus() == 0
		}
		c.Log.Warnw("something went wrong with 'which', trying native command")
		return c.checkPackageManagerWithNativeCommand(manager)
	}
	waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
	return waitStatus.ExitStatus() == 0
}

func (c *cliState) checkPackageManagerWithNativeCommand(manager string) bool {
	var (
		cmd    = exec.Command("command", "-v", manager)
		_, err = cmd.CombinedOutput()
	)
	if err != nil {
		c.Log.Debugw("error trying to check package-manager",
			"cmd", "command",
			"package-manager", manager,
			"error", err,
		)
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			return waitStatus.ExitStatus() == 0
		}
		return false
	}
	waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
	return waitStatus.ExitStatus() == 0
}

func removeEpochFromPkgVersion(pkgVer string) string {
	if strings.Contains(pkgVer, ":") {
		pkgVerSplit := strings.Split(pkgVer, ":")
		if len(pkgVerSplit) == 2 {
			return pkgVerSplit[1]
		}
	}

	return pkgVer
}

// split the provided package_manifest into chucks, if the manifest
// is smaller than the provided chunk size, it will return the manifest
// as an array without modifications
func splitPackageManifest(manifest *api.VulnerabilitiesPackageManifest, chunks int) []*api.VulnerabilitiesPackageManifest {
	if len(manifest.OsPkgInfoList) <= chunks {
		return []*api.VulnerabilitiesPackageManifest{manifest}
	}

	var batches []*api.VulnerabilitiesPackageManifest
	for i := 0; i < len(manifest.OsPkgInfoList); i += chunks {
		batch := manifest.OsPkgInfoList[i:min(i+chunks, len(manifest.OsPkgInfoList))]
		cli.Log.Infow("manifest batch", "total_packages", len(batch))
		batches = append(batches, &api.VulnerabilitiesPackageManifest{OsPkgInfoList: batch})
	}
	return batches
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// fan-out a number of package manifests into multiple requests all at once
func fanOutHostScans(manifests ...*api.VulnerabilitiesPackageManifest) (api.VulnerabilitySoftwarePackagesResponse, error) {
	var (
		resCh    = make(chan api.VulnerabilitySoftwarePackagesResponse)
		errCh    = make(chan error)
		workers  = len(manifests)
		fanInRes = api.VulnerabilitySoftwarePackagesResponse{}
	)

	// disallow more than 10 workers which are 10 calls all at once,
	// the API has a rate-limit of 10 calls per hour, per access key
	if workers > 10 {
		return fanInRes, errors.New("limit of packages exceeded")
	}

	var (
		err   error
		start = time.Now()
	)
	defer func() {
		cli.Event.DurationMs = time.Since(start).Milliseconds()
		// avoid duplicating events
		if err == nil {
			cli.SendHoneyvent()
		}
	}()

	// ensure that the api client has a valid token
	// before creating workers
	if !cli.LwApi.ValidAuth() {
		_, err = cli.LwApi.GenerateToken()
		if err != nil {
			return fanInRes, err
		}
	}

	// for every manifest, create a new worker, that is, spawn
	// a new goroutine that will send the manifest to scan
	for n, m := range manifests {
		if m == nil {
			workers--
			continue
		}
		cli.Log.Infow("spawn worker", "number", n+1)
		go cli.triggerHostVulnScan(m, resCh, errCh)
	}

	cli.Event.AddFeatureField("workers", workers)

	// lock the main process and read both, the error and response
	// channels, if we receive at least one error, we will stop
	// processing and bubble up the error to the caller
	for processed := 0; processed < workers; processed++ {
		select {
		case err = <-errCh:
			// end processing as soon as we receive the first error
			return fanInRes, err
		case res := <-resCh:
			// processing scan
			cli.Log.Infow("processing worker response", "n", processed+1)
			cli.Event.AddFeatureField(fmt.Sprintf("worker%d_total_vulns", processed), len(res.Data))
			mergeHostVulnScanPkgManifestResponses(&fanInRes, &res)
		}
	}

	return fanInRes, nil
}

func mergeHostVulnScanPkgManifestResponses(to, from *api.VulnerabilitySoftwarePackagesResponse) {
	// append vulnerabilities from -> to
	to.Data = append(to.Data, from.Data...)

	// Todo(v2): this is no longer relevant
	// requests should always return an ok state
	//to.Ok = from.Ok

	// store the message from the response only if it is NOT empty
	// and it is different from the previous response (to)
	//if to.Message == "" {
	//	to.Message = from.Message
	//	return
	//}

	// concatenate messages "to,from" response only if they
	// are NOT empty and they are different from each other
	//if from.Message != "" && from.Message != to.Message {
	//	to.Message = fmt.Sprintf("%s,%s", to.Message, from.Message)
	//}
}

func (c *cliState) triggerHostVulnScan(manifest *api.VulnerabilitiesPackageManifest,
	resCh chan<- api.VulnerabilitySoftwarePackagesResponse,
	errCh chan<- error,
) {
	response, err := c.LwApi.V2.Vulnerabilities.SoftwarePackages.Scan(*manifest)
	if err != nil {
		errCh <- err
		return
	}
	resCh <- response
}
