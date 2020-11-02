// +build !integration

/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package properties_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/properties-file-image-relocation/pkg/internal/properties"
)

func TestImages(t *testing.T) {
	data, err := properties.Images(mustReadFile(t, "./test_data/props"))
	require.NoError(t, err)
	require.ElementsMatch(t, []string{
		"springcloudstream/cassandra-sink-rabbit:2.1.2.RELEASE",
		"springcloudstream/counter-sink-rabbit:2.1.2.RELEASE",
		"springcloudstream/file-sink-rabbit:2.1.2.RELEASE",
		"springcloudstream/ftp-sink-rabbit:2.1.2.RELEASE",
	}, data)
}

func TestInvalidImage(t *testing.T) {
	// Image references should include all text to the end of the line,
	// even though that will produce an invalid image reference.
	// This gives early warning of junk on the end of lines.
	data, err := properties.Images(mustReadFile(t, "./test_data/props.invalidref"))
	require.NoError(t, err)
	require.ElementsMatch(t, []string{
		"imagename extra-text",
	}, data)
}

func mustReadFile(t *testing.T, file string) []byte {
	data, err := ioutil.ReadFile(file)
	require.NoError(t, err)
	return data
}
