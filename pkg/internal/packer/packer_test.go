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
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/properties-file-image-relocation/pkg/internal/packer"
)

func TestPackUnpack(t *testing.T) {
	work, err := ioutil.TempDir("", "prel-packer-test")
	require.NoError(t, err)

	archivePath := filepath.Join(work, "test.tgz")

	err = packer.Pack("./test_data/props", archivePath)
	require.NoError(t, err)

	unpacked, propsFile, err := packer.Unpack(archivePath)
	require.NoError(t, err)
	defer os.RemoveAll(unpacked)

	data, err := ioutil.ReadFile(propsFile)
	require.NoError(t, err)

	require.Equal(t, fmt.Sprintf("# properties file with no images so it can be used in a unit test%sa=b", endOfLine()), string(data))
}

func endOfLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
