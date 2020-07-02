// +build !integration

/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package packer_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/pivotal/scdf-k8s-prel/pkg/internal/packer"
	"github.com/stretchr/testify/require"
)

func TestPackUnpack(t *testing.T) {
	work, err := ioutil.TempDir("", "prel-packer-test")
	require.NoError(t, err)

	archivePath := filepath.Join(work, "test.tgz")

	err = packer.Pack("./test_data/props", archivePath)
	require.NoError(t, err)

	unpacked, err := packer.Unpack(archivePath)
	require.NoError(t, err)

	data, err := ioutil.ReadFile(filepath.Join(unpacked, "props"))
	require.NoError(t, err)

	require.Equal(t, fmt.Sprintf("# properties file with no images so it can be used in a unit test%sa=b", endOfLine()), string(data))
}

func endOfLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}