/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package properties

import (
	"strings"

	mprops "github.com/magiconair/properties"
	"github.com/pivotal/scdf-k8s-prel/pkg/internal/config"
)

// Images returns the image references in the given properties
func Images(data []byte) ([]string, error) {
	props := mprops.MustLoadString(string(data))
	imageRefs := map[string]struct{}{}        // use a map to remove duplicates
	props.FilterFunc(func(_, v string) bool { // use FilterFunc to iterate over the values in props
		if found, imageRef, _ := detectImage(v); found {
			imageRefs[imageRef] = struct{}{}
		}
		return false // any return value will do
	})

	result := []string{}
	for ref := range imageRefs {
		result = append(result, ref)
	}
	return result, nil
}

func detectImage(value string) (bool, string, string) {
	imageValuePrefixes := config.PropertyValueImagePrefixes()
	// if value has a prefix which means it is an image reference, record the length of the longest
	// matching prefix. This ensures, for example, that "docker:" is not taken to be a prefix of
	// "docker://imageref" when the prefix is actually "docker://".
	longestPrefix := ""
	for _, prefix := range imageValuePrefixes {
		if strings.HasPrefix(value, prefix) {
			if len(prefix) > len(longestPrefix) {
				longestPrefix = prefix
			}
		}
	}
	// if there was a match, strip off the longest matching prefix and return the image
	// reference and the prefix
	if longestPrefix != "" {
		return true, value[len(longestPrefix):], longestPrefix
	}

	return false, "", ""
}
