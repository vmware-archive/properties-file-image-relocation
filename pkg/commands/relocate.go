/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"log"

	"github.com/pivotal/scdf-k8s-prel/pkg/internal/relocator"
	"github.com/spf13/cobra"
)

// NewRelocateCommand creates a relocate command.
func NewRelocateCommand() *cobra.Command {
	var repositoryPrefix, outputPath string
	cmd := &cobra.Command{
		Use:   "relocate ZIPPED-ARCHIVE-PATH",
		Short: "Relocate the properties file and images in a zipped archive",
		Long: `Relocate any docker and OCI images in a zipped archive (which was created using the package
subcommand), push the images to a registry, and output a relocated properties file.

The --repository-prefix flag determines the repositories for the relocated images. Each image is assigned
a name starting with the given prefix and pushed to the repository. For example, if the repository prefix
is example.com/user, each image is relocated to a name starting with example.com/user/ and pushed to a
repository hosted by example.com.

The --output flag specifies the path where the output UTF-8 encoded relocated properties file will be writted.
This properties file has the same properties as the properties file in the zipped archive except that image
references are replaced by their relocated counterparts.
`,
		Args: cobra.ExactArgs(1), // the zipped archive
		Run: func(_ *cobra.Command, args []string) {
			if err := relocator.Relocate(args[0], repositoryPrefix, outputPath); err != nil {
				log.Fatalf("relocate command failed: %v", err)
			}
		},
	}

	cmd.Flags().StringVarP(&repositoryPrefix, "repository-prefix", "p", "", "prefix for relocated image names")
	cmd.MarkFlagRequired("repository-prefix")

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "path for output properties file")
	cmd.MarkFlagRequired("output")

	return cmd
}
