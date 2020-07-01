// +build !integration

/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	ver := CliVersion()
	require.Equal(t, "unknown (unknown sha)", ver)
}
