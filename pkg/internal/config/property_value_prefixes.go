/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package config

// PropertyValueImagePrefixes returns the property value prefixes which denote images
func PropertyValueImagePrefixes() []string {
	return []string{"docker:", "docker://"}
}
