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

	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/properties-file-image-relocation/pkg/internal/packer"
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
	if output, err := exec.Command(prelBin, "package", "--file", "./test_data/props", "--archive", archivePath).CombinedOutput(); err != nil {
		require.NoError(t, err, string(output))
	}

	unpacked, propsPath, err := packer.Unpack(archivePath)
	require.NoError(t, err)
	defer os.RemoveAll(unpacked)

	data, err := ioutil.ReadFile(propsPath)
	require.NoError(t, err)

	expected := "# integration test with small images\na=b\nalpine=docker:amd64/alpine:latest\nc=d\nhello:docker://amd64/hello-world:latest\ne=f"
	expected = strings.ReplaceAll(expected, "\n", endOfLine())
	require.Equal(t, expected, string(data))
}

func TestRelocate(t *testing.T) {
	// If no registry is available, skip the test.
	registry, registryPresent := os.LookupEnv("REGISTRY")
	if !registryPresent {
		t.Log("Skipping relocation integration test since REGISTRY environment variable is not set")
		return
	}

	work, err := ioutil.TempDir("", "prel-integration-test")
	require.NoError(t, err)

	archivePath := filepath.Join(work, "test.tgz")

	prelBin := os.Getenv("PREL_EXECUTABLE")
	if output, err := exec.Command(prelBin, "package", "--file", "./test_data/props", "--archive", archivePath).CombinedOutput(); err != nil {
		require.NoError(t, err, string(output))
	}

	relocatedPropertiesPath := filepath.Join(work, "props.relocated")

	if output, err := exec.Command(prelBin, "relocate", "--archive", archivePath, "--repository-prefix", registry+"/user", "--output", relocatedPropertiesPath).CombinedOutput(); err != nil {
		require.NoError(t, err, string(output))
	}

	data, err := ioutil.ReadFile(relocatedPropertiesPath)
	require.NoError(t, err)

	expected := fmt.Sprintf("# integration test with small images\na = b\nalpine = docker:%s/user/amd64-alpine-2ce5bd4449b9abea56a42cfb3d073c8e:latest\nc = d\n"+
		"hello = docker://%s/user/amd64-hello-world-c2d67489afc2278fc40cab6d4d34b521:latest\ne = f\n", registry, registry)
	expected = strings.ReplaceAll(expected, "\n", endOfLine())
	require.Equal(t, expected, string(data))
}

func endOfLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
