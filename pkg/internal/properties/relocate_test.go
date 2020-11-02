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

	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/properties-file-image-relocation/pkg/internal/properties"
)

func TestRelocate(t *testing.T) {
	mapping := map[string]string{
		"springcloudstream/cassandra-sink-rabbit:2.1.2.RELEASE": "xspringcloudstream/cassandra-sink-rabbit:2.1.2.RELEASE",
		"springcloudstream/counter-sink-rabbit:2.1.2.RELEASE":   "xspringcloudstream/counter-sink-rabbit:2.1.2.RELEASE",
		"springcloudstream/file-sink-rabbit:2.1.2.RELEASE":      "xspringcloudstream/file-sink-rabbit:2.1.2.RELEASE",
		"springcloudstream/ftp-sink-rabbit:2.1.2.RELEASE":       "xspringcloudstream/ftp-sink-rabbit:2.1.2.RELEASE",
	}
	data, err := properties.Relocate("./test_data/props", mapping)
	require.NoError(t, err)
	require.Equal(t, `# cassandra comment
sink.cassandra = docker:xspringcloudstream/cassandra-sink-rabbit:2.1.2.RELEASE
sink.cassandra.metadata = maven://org.springframework.cloud.stream.app:cassandra-sink-rabbit:jar:metadata:2.1.2.RELEASE

# counter
# comments
sink.counter = docker:xspringcloudstream/counter-sink-rabbit:2.1.2.RELEASE
sink.counter.metadata = maven://org.springframework.cloud.stream.app:counter-sink-rabbit:jar:metadata:2.1.2.RELEASE
sink.file = docker://xspringcloudstream/file-sink-rabbit:2.1.2.RELEASE
sink.file.metadata = maven://org.springframework.cloud.stream.app:file-sink-rabbit:jar:metadata:2.1.2.RELEASE
sink.ftp = docker:xspringcloudstream/ftp-sink-rabbit:2.1.2.RELEASE
sink.ftp.metadata = maven://org.springframework.cloud.stream.app:ftp-sink-rabbit:jar:metadata:2.1.2.RELEASE
another.cassandra = docker://xspringcloudstream/cassandra-sink-rabbit:2.1.2.RELEASE
`, string(data))
}

func TestRelocateUnmapped(t *testing.T) {
	mapping := map[string]string{}

	_, err := properties.Relocate("./test_data/props.single", mapping)
	require.Equal(t, "error(s) relocating properties file: image in 'single=docker:imagename' was not relocated", err.Error())
}

func TestRelocateFileNotFound(t *testing.T) {
	_, err := properties.Relocate("./test_data/no-such", nil)
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
}
