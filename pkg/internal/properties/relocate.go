/*
 * Copyright 2020 VMware, Inc.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package properties

import (
	"bytes"
	"fmt"

	mprops "github.com/magiconair/properties"
	"github.com/pivotal/go-ape/pkg/furl"
)

// Relocate relocates the properties file at the given path according to the given mapping
// of image references to relocated image references and returns the relocated properties
// as data ready to write to a properties file.
func Relocate(propsFile string, mapping map[string]string) ([]byte, error) {
	data, err := furl.Read(propsFile, "")
	if err != nil {
		return nil, err
	}

	props := mprops.MustLoadString(string(data))
	errs := []error{}
	props.FilterFunc(func(k, v string) bool { // use FilterFunc to iterate over the values in props
		if found, imageRef, prefix := detectImage(v); found {
			relocated, ok := mapping[imageRef]
			var newV string
			if ok {
				newV = prefix + relocated
			} else {
				errs = append(errs, fmt.Errorf("image in '%s=%s' was not relocated", k, v))
				return false // any value will do
			}
			_, _, err := props.Set(k, newV)
			if err != nil {
				errs = append(errs, err)
			}
		} else {
			_, _, err := props.Set(k, v)
			if err != nil {
				errs = append(errs, err)
			}
		}

		return false // any return value will do
	})

	if len(errs) > 0 {
		message := ""
		for i, err := range errs {
			if i > 0 {
				message += "; "
			}
			message += err.Error()
		}
		return nil, fmt.Errorf("error(s) relocating properties file: %s", message)
	}

	var buf bytes.Buffer
	if _, err := props.WriteComment(&buf, "# ", mprops.UTF8); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
