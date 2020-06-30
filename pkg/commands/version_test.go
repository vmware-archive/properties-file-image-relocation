/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands_test

import (
	"testing"

	"github.com/pivotal/scdf-k8s-prel/pkg/commands"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	ver := commands.CliVersion()
	require.Equal(t, "unknown (unknown sha)", ver)
}
