// +build integration

/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/pivotal/scdf-k8s-prel/pkg/internal/packer"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	if err := os.Chdir("../.."); err != nil {
		fmt.Printf("change directory failed: %v\n", err)
		os.Exit(1)
	}
	if err := exec.Command("make", "prel").Run(); err != nil {
		fmt.Printf("make prel failed: %v\n", err)
		os.Exit(1)
	}
	root, err := os.Getwd()
	if err != nil {
		fmt.Printf("get working directory failed: %v\n", err)
		os.Exit(1)
	}
	bin := filepath.Join(root, "bin", "prel")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	os.Setenv("PREL_EXECUTABLE", bin)
	if err := os.Chdir("pkg/integration"); err != nil {
		fmt.Printf("change directory failed: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestPackage(t *testing.T) {
	work, err := ioutil.TempDir("", "prel-integration-test")
	require.NoError(t, err)

	archivePath := filepath.Join(work, "test.tgz")

	prelBin := os.Getenv("PREL_EXECUTABLE")
	if output, err := exec.Command(prelBin, "package", "./test_data/props", "--archive", archivePath).CombinedOutput(); err != nil {
		require.NoError(t, err, string(output))
	}

	unpacked, err := packer.Unpack(archivePath)
	require.NoError(t, err)

	data, err := ioutil.ReadFile(filepath.Join(unpacked, "props"))
	require.NoError(t, err)

	expected := "# integration test with small images\na=b\nalpine=docker:amd64/alpine:latest\nc=d\ndocker:oddlynamed=docker://amd64/hello-world:latest\ne=f"
	expected = strings.ReplaceAll(expected, "\n", endOfLine())
	require.Equal(t, expected, string(data))
}

func endOfLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
