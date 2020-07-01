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

func TestImageDetection(t *testing.T) {
	data, err := properties.Images("./test_data/props")
	require.NoError(t, err)
	require.ElementsMatch(t, []string{
		"springcloudstream/cassandra-sink-rabbit:2.1.2.RELEASE",
		"springcloudstream/counter-sink-rabbit:2.1.2.RELEASE",
		"springcloudstream/file-sink-rabbit:2.1.2.RELEASE",
		"springcloudstream/ftp-sink-rabbit:2.1.2.RELEASE",
	}, data)
}

func TestImageDetectionFileNotFound(t *testing.T) {
	_, err := properties.Images("./test_data/no-such")
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
}
