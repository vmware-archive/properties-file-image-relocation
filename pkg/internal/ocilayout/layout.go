/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package ocilayout

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/registry/ggcr"
)

// StoreImages stores the images with the given references in the given directory
func StoreImages(archiveDir string, imageRefs []string) error {
	layoutDir := filepath.Join(archiveDir, "layout")
	if err := os.MkdirAll(layoutDir, 0755); err != nil {
		return err
	}

	layout, err := ggcr.NewRegistryClient().NewLayout(layoutDir)
	if err != nil {
		return err
	}

	for _, ref := range imageRefs {
		nm, err := image.NewName(ref)
		if err != nil {
			return err
		}

		fmt.Printf("Downloading %s\n", ref)
		if _, err := layout.Add(nm); err != nil {
			return err
		}
	}

	return nil
}
