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
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"syscall"

	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/api"
)

var SupportedPackageManagers = []string{"dpkg-query", "rpm"} // @afiune can we support ym and apk?

type OS struct {
	Name    string
	Version string
}

var (
	osReleaseFile = "/etc/os-release"
	rexNameFromID = regexp.MustCompile(`^ID=(.*)$`)
	rexVersionID  = regexp.MustCompile(`^VERSION_ID=(.*)$`)
)

func (c *cliState) GeneratePackageManifest() (*api.PackageManifest, error) {
	manifest := new(api.PackageManifest)
	osInfo, err := cli.GetOSInfo()
	if err != nil {
		return manifest, err
	}
	manager, err := cli.DetectPackageManager()
	if err != nil {
		return manifest, err
	}

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
			"this is most likely a mistake on us, please report it to support.lacework.com.",
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
			api.OsPkgInfo{
				Os:     osInfo.Name,
				OsVer:  osInfo.Version,
				Pkg:    pkgDetail[0],
				PkgVer: pkgDetail[1],
			},
		)
	}

	c.Log.Debugw("package-manifest", "raw", manifest)
	return manifest, nil
}

func (c *cliState) GetOSInfo() (*OS, error) {
	osInfo := new(OS)

	c.Log.Debugw("detecting operating system information",
		"os", runtime.GOOS,
		"arch", runtime.GOARCH,
	)

	f, err := os.Open(osReleaseFile)
	if err != nil {
		msg := `unsupported platform

For more information about supported platforms, visit:
    https://support.lacework.com/hc/en-us/articles/360049666194-Host-Vulnerability-Assessment-Overview`
		return osInfo, errors.New(msg)
	}
	defer f.Close()

	c.Log.Debugw("parsing os release file", "file", osReleaseFile)
	s := bufio.NewScanner(f)
	for s.Scan() {
		if m := rexNameFromID.FindStringSubmatch(s.Text()); m != nil {
			osInfo.Name = strings.Trim(m[1], `"`)
		} else if m := rexVersionID.FindStringSubmatch(s.Text()); m != nil {
			osInfo.Version = strings.Trim(m[1], `"`)
		}
	}

	return osInfo, nil
}

func (c *cliState) DetectPackageManager() (string, error) {
	c.Log.Debugw("detecting package-manager")

	for _, manager := range SupportedPackageManagers {
		if cli.checkPackageManager(manager) {
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
		return cli.checkPackageManagerWithNativeCommand(manager)
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
