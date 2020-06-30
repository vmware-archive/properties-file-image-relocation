/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import "fmt"

var (
	cliVersion  = "unknown"
	cliGitsha   = "unknown sha"
	cliGitdirty = ""
)

// CliVersion returns a version string based on values supplied at build time
func CliVersion() string {
	var version string
	if cliGitdirty == "" {
		version = fmt.Sprintf("%s (%s)", cliVersion, cliGitsha)
	} else {
		version = fmt.Sprintf("%s (%s, with local modifications)", cliVersion, cliGitsha)
	}
	return version
}
