/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/properties-file-image-relocation/pkg/internal/packer"
)

// NewPackageCommand creates a package command.
func NewPackageCommand() *cobra.Command {
	var (
		archivePath    string
		propertiesFile string
	)
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Package a properties file and its images in a zipped archive",
		Long:  `Package a UTF-8 encoded properties file and its docker and OCI images in a zipped archive`,
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, args []string) {
			if err := packer.Pack(propertiesFile, archivePath); err != nil {
				log.Fatalf("package command failed: %v", err)
			}
		},
	}

	cmd.Flags().StringVarP(&propertiesFile, "file", "f", "", "file path or URL of properties file, or - to read from standard input")
	cmd.MarkFlagRequired("file")

	cmd.Flags().StringVarP(&archivePath, "archive", "a", "", "file path of the archive to be created")
	cmd.MarkFlagRequired("archive")
	return cmd
}
