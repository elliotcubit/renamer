/*
   Copyright 2020 Ori Seri

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package regexps

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	requiredOption = "required"
	existsOption   = "exists"
)

// matchGroupMap converts the match string array into a map of group keys to group values
func (r *Regexp[T]) matchGroupMap(match []string) map[string]string {
	ret := make(map[string]string)
	for i, name := range r.matcher.SubexpNames() {
		if i != 0 && name != "" {
			ret[name] = match[i]
		}
	}
	return ret
}

// groupAndOption returns the requested regexps and its options split by ','
func groupAndOption(fieldType reflect.StructField) (group string, option []string) {
	regroupKey := fieldType.Tag.Get("regexps")
	if regroupKey == "" {
		return "", nil
	}
	split := strings.Split(regroupKey, ",")
	if len(split) == 1 {
		return strings.TrimSpace(split[0]), nil
	}
	var options []string
	for _, opt := range split[1:] {
		options = append(options, strings.TrimSpace(opt))
	}
	return strings.TrimSpace(split[0]), options
}

// validateStruct checks that all fields in the given struct are valid
func (r *Regexp[T]) validateStruct(targetRef reflect.Value) error {
	targetType := targetRef.Type()
	for i := 0; i < targetType.NumField(); i++ {
		fieldRef := targetRef.Field(i)
		if !fieldRef.CanSet() {
			continue
		}
		if err := r.validateField(targetType.Field(i), fieldRef); err != nil {
			return err
		}
	}
	return nil
}

// validateField ensures that the given field is present if required, or (TODO) has a default value
func (r *Regexp[T]) validateField(fieldType reflect.StructField, fieldRef reflect.Value) error {
	fieldRefType := fieldType.Type
	ptr := false
	if fieldRefType.Kind() == reflect.Ptr {
		ptr = true
		fieldRefType = fieldType.Type.Elem()
	}
	if fieldRefType.Kind() == reflect.Struct {
		if ptr {
			if fieldRef.IsNil() {
				return fmt.Errorf("can't set value to nil pointer in struct field: %s", fieldType.Name)
			}
			fieldRef = fieldRef.Elem()
		}
		return r.validateStruct(fieldRef)
	}
	key, opts := groupAndOption(fieldType)
	if key == "" {
		return nil
	}
	if slices.Contains(opts, requiredOption) {
		if !slices.Contains(r.matcher.SubexpNames(), key) {
			if _, ok := r.defaults[key]; !ok {
				return fmt.Errorf("missing required field %q", key)
			}
		}
	}
	return nil
}

// setField getting a single struct field and matching groups map and set the field value to its matching group value tag
// after parsing it to match the field type
func (r *Regexp[T]) setField(fieldType reflect.StructField, fieldRef reflect.Value, matchGroup map[string]string) error {
	fieldRefType := fieldType.Type
	ptr := false
	if fieldRefType.Kind() == reflect.Ptr {
		fieldRefType = fieldType.Type.Elem()
	}

	if fieldRefType.Kind() == reflect.Struct {
		return r.fillTarget(matchGroup, fieldRef)
	}

	key, _ := groupAndOption(fieldType)
	if key == "" {
		return nil
	}

	if ptr {
		fieldRef = fieldRef.Elem()
	}

	matchedVal, ok := matchGroup[key]
	if !ok {
		matchedVal = r.defaults[key]
	}

	parsedFunc := getParsingFunc(fieldRefType)
	if parsedFunc == nil {
		return &TypeNotParsableError{fieldRefType}
	}

	parsed, err := parsedFunc(matchedVal, fieldRefType)
	if err != nil {
		return &ParseError{group: key, err: err}
	}

	fieldRef.Set(parsed)

	return nil
}

func (r *Regexp[T]) fillTarget(matchGroup map[string]string, targetRef reflect.Value) error {
	targetType := targetRef.Type()
	for i := 0; i < targetType.NumField(); i++ {
		fieldRef := targetRef.Field(i)
		if !fieldRef.CanSet() {
			continue
		}

		if err := r.setField(targetType.Field(i), fieldRef, matchGroup); err != nil {
			return err
		}
	}

	return nil
}

// validateTarget checks that given interface is a pointer of struct
func validateTarget(target interface{}) (reflect.Value, error) {
	targetPtr := reflect.ValueOf(target)
	if targetPtr.Kind() != reflect.Ptr {
		return reflect.Value{}, &NotStructPtrError{}
	}
	return targetPtr.Elem(), nil
}
