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
	"github.com/pivotal/image-relocation/pkg/pathmapping"
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

	for i, ref := range imageRefs {
		nm, err := image.NewName(ref)
		if err != nil {
			return err
		}

		fmt.Printf("Downloading image (%d of %d) %s\n", i+1, len(imageRefs), ref)
		if _, err := layout.Add(nm); err != nil {
			return err
		}
	}

	return nil
}

// RelocateImages relocates the given image references using the given repository prefix and pushes
// them from the OCI image layout in the given archive directory to their relocated repositories. It
// returns a mapping from original to relocated image reference.
func RelocateImages(archiveDir string, imageRefs []string, prefix string) (map[string]string, error) {
	layoutDir := filepath.Join(archiveDir, "layout")

	l, err := ggcr.NewRegistryClient().ReadLayout(layoutDir)
	if err != nil {
		return nil, err
	}

	mapping := map[string]string{}

	for i, imageRef := range imageRefs {
		imageName, err := image.NewName(imageRef)
		if err != nil {
			return nil, err
		}
		imageDigest, err := l.Find(imageName)
		if err != nil {
			return nil, fmt.Errorf("image reference %s not found in archive: %w", imageName.String(), err)
		}

		newImageName, err := pathmapping.FlattenRepoPathPreserveTagDigest(prefix, imageName)
		if err != nil {
			return nil, fmt.Errorf("image reference %s could not be relocated: %w", imageName.String(), err)
		}

		fmt.Printf("Relocating image (%d of %d) %s to %s \n", i+1, len(imageRefs), imageRef, newImageName.String())
		if err := l.Push(imageDigest, newImageName); err != nil {
			return nil, fmt.Errorf("failed to push to repository %s: %w", newImageName.String(), err)
		}
		mapping[imageRef] = newImageName.String()
	}

	return mapping, nil
}
