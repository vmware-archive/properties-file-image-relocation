/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"github.com/spf13/cobra"
)

// NewRootCommand creates a root for the command hierarchy.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "prel",
		Short:        "prel packages and relocates properties files",
		SilenceUsage: true,
		Run:          func(cmd *cobra.Command, _ []string) { cmd.Usage() },
	}

	cmd.Version = CliVersion()
	cmd.Flags().Bool("version", false, "display command version")

	cmd.AddCommand(NewPackageCommand())

	return cmd
}

// Execute executes the root command for prel.
func Execute() error {
	return NewRootCommand().Execute()
}
