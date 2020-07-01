/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package packer

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/moby/moby/pkg/archive"
	"github.com/pivotal/go-ape/pkg/furl"
	"github.com/pivotal/scdf-k8s-prel/pkg/internal/ocilayout"
	"github.com/pivotal/scdf-k8s-prel/pkg/internal/properties"
)

// Pack packs the given properties file and its images in a tgz archive with the given file path
func Pack(props, archivePath string) error {
	archiveDir, err := ioutil.TempDir("", "prel-packer")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(archiveDir)

	propsData, err := furl.Read(props, "")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(archiveDir, "props"), propsData, 0666); err != nil {
		return err
	}

	imageRefs, err := properties.Images(props)
	if err != nil {
		return err
	}

	if err := ocilayout.StoreImages(archiveDir, imageRefs); err != nil {
		return err
	}

	writer, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("Error creating archive file: %s", err)
	}
	defer writer.Close()

	tarOptions := &archive.TarOptions{
		Compression:      archive.Gzip,
		IncludeFiles:     []string{"."},
		IncludeSourceDir: true,
	}
	rc, err := archive.TarWithOptions(archiveDir, tarOptions)
	if err != nil {
		return err
	}
	defer rc.Close()

	_, err = io.Copy(writer, rc)
	return err
}
