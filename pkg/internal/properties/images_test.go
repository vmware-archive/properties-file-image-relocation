// +build !integration

/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package properties_test

import (
	"os"
	"testing"

	"github.com/pivotal/scdf-k8s-prel/pkg/internal/properties"
	"github.com/stretchr/testify/require"
)

func TestImages(t *testing.T) {
	data, err := properties.Images("./test_data/props")
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
	data, err := properties.Images("./test_data/props.invalidref")
	require.NoError(t, err)
	require.ElementsMatch(t, []string{
		"imagename extra-text",
	}, data)
}

func TestImagesFileNotFound(t *testing.T) {
	_, err := properties.Images("./test_data/no-such")
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
}
