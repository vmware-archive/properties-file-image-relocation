/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package packer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/moby/moby/pkg/archive"
)

// Unpack decompresses the given tgz archive (created using Pack) to a temporary directory and
// returns the directory and the path of the properties file in the directory.
// It is the caller's responsibility to delete the temporary directory
func Unpack(archiveDir string) (string, string, error) {
	unpackDir, err := ioutil.TempDir("", "prel-packer")
	if err != nil {
		return "", "", err
	}

	reader, err := os.Open(archiveDir)
	if err != nil {
		return "", "", err
	}
	defer reader.Close()

	tarOptions := &archive.TarOptions{
		Compression:      archive.Gzip,
		IncludeFiles:     []string{"."},
		IncludeSourceDir: true,
		// Issue #416
		NoLchown: true,
	}
	if err := archive.Untar(reader, unpackDir, tarOptions); err != nil {
		return "", "", fmt.Errorf("untar failed: %s", err)
	}

	return unpackDir, filepath.Join(unpackDir, propertiesFilePath), nil
}
