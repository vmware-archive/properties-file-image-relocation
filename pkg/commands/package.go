/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"log"

	"github.com/pivotal/scdf-k8s-prel/pkg/internal/packer"
	"github.com/spf13/cobra"
)

// NewPackageCommand creates a package command.
func NewPackageCommand() *cobra.Command {
	var archivePath string
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Package a properties file and its images into a zipped archive",
		Args:  cobra.ExactArgs(1), // the properties file
		Run: func(_ *cobra.Command, args []string) {
			if err := packer.Pack(args[0], archivePath); err != nil {
				log.Fatalf("package command failed: %v", err)
			}
		},
	}
	cmd.Flags().StringVarP(&archivePath, "archive", "a", "", "file path of the archive to be created")
	cmd.MarkFlagRequired("archive")
	return cmd
}
