/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package relocator

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/vmware-tanzu/properties-file-image-relocation/pkg/internal/ocilayout"
	"github.com/vmware-tanzu/properties-file-image-relocation/pkg/internal/packer"
	"github.com/vmware-tanzu/properties-file-image-relocation/pkg/internal/properties"
)

// Relocate relocates the images in the zipped archive at the given path by
// applying the given repository prefix. It creates a relocated properties
// file at the output path.
func Relocate(archivePath, repositoryPrefix, outputPath string) error {
	unpacked, propsFile, err := packer.Unpack(archivePath)
	if err != nil {
		return err
	}
	defer os.RemoveAll(unpacked)

	propsData, err := ioutil.ReadFile(propsFile)
	if err != nil {
		return err
	}

	imageRefs, err := properties.Images(propsData)
	if err != nil {
		return err
	}

	mapping, err := ocilayout.RelocateImages(unpacked, imageRefs, repositoryPrefix)
	if err != nil {
		return err
	}

	relocatedProperties, err := properties.Relocate(propsFile, mapping)
	if err != nil {
		return err
	}

	fmt.Printf("Creating relocated properties file %s\n", outputPath)
	return ioutil.WriteFile(outputPath, relocatedProperties, 0666)
}
