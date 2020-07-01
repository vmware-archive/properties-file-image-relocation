/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package properties

import (
	"strings"
	"unicode"

	"github.com/magiconair/properties"
	"github.com/pivotal/go-ape/pkg/furl"
	"github.com/pivotal/scdf-k8s-prel/pkg/internal/config"
)

// Images returns the image references in the given properties file, which may be specified by file path or URL
func Images(propsFile string) ([]string, error) {
	data, err := furl.Read(propsFile, "")
	if err != nil {
		return []string{}, err
	}

	props := properties.MustLoadString(string(data))
	imageValuePrefixes := config.PropertyValueImagePrefixes()
	imageRefs := map[string]struct{}{}        // use a map to remove duplicates
	props.FilterFunc(func(_, v string) bool { // use FilterFunc to iterate over the values in props
		// if v has a prefix which means it is an image reference, record the length of the longest
		// matching prefix. This ensures, for example, that "docker:" is not taken to be a prefix of
		// "docker://imageref" when the prefix is actually "docker://".
		longestPrefixLength := 0
		for _, prefix := range imageValuePrefixes {
			if strings.HasPrefix(v, prefix) {
				if len(prefix) > longestPrefixLength {
					longestPrefixLength = len(prefix)
				}
			}
		}
		// if there was a match, strip off the longest matching prefix and record the image
		// reference, stopping at the first whitespace character, if there is one, or at
		// the end of the string otherwise.
		if longestPrefixLength > 0 {
			suffix := v[longestPrefixLength:]
			imageRef := ""
			for _, s := range suffix {
				if unicode.IsSpace(s) {
					break
				}
				imageRef += string(s)
			}
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
