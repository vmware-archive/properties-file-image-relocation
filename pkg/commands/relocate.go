/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/properties-file-image-relocation/pkg/internal/relocator"
)

// NewRelocateCommand creates a relocate command.
func NewRelocateCommand() *cobra.Command {
	var (
		archivePath      string
		repositoryPrefix string
		outputPath       string
	)
	cmd := &cobra.Command{
		Use:   "relocate",
		Short: "Relocate the properties file and images in a zipped archive",
		Long: `Relocate any docker and OCI images in a zipped archive (which was created using the package
subcommand), push the images to a registry, and output a relocated properties file.

The --archive flag specifies the file path of the zipped archive.

The --repository-prefix flag determines the repositories for the relocated images. Each image is assigned
a name starting with the given prefix and pushed to the repository. For example, if the repository prefix
is example.com/user, each image is relocated to a name starting with example.com/user/ and pushed to a
repository hosted by example.com.

The --output flag specifies the file path where the output UTF-8 encoded relocated properties file will be
written. This properties file has the same properties as the properties file in the zipped archive except
that image references are replaced by their relocated counterparts.
`,
		Args: cobra.NoArgs,
		Run: func(_ *cobra.Command, args []string) {
			if err := relocator.Relocate(archivePath, repositoryPrefix, outputPath); err != nil {
				log.Fatalf("relocate command failed: %v", err)
			}
		},
	}

	cmd.Flags().StringVarP(&archivePath, "archive", "a", "", "file path of zipped archive")
	cmd.MarkFlagRequired("archive")

	cmd.Flags().StringVarP(&repositoryPrefix, "repository-prefix", "p", "", "prefix for relocated image names")
	cmd.MarkFlagRequired("repository-prefix")

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "file path for output properties file")
	cmd.MarkFlagRequired("output")

	return cmd
}
