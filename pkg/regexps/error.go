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
)

// CompileError returned on regex compilation error
type CompileError struct{ err error }

func (c *CompileError) Error() string {
	return fmt.Sprintf("compilation error: %v", c.err)
}

// NoMatchFoundError indicates no regex matches for given string
type NoMatchFoundError struct{}

func (n *NoMatchFoundError) Error() string {
	return "no match found for given string"
}

// NotStructPtrError returned when given target is not a truct pointer
type NotStructPtrError struct{}

func (n *NotStructPtrError) Error() string {
	return "expected struct pointer"
}

// UnknownGroupError returned when given regex group tag isn't exists in compiled regex groups
type UnknownGroupError struct{ group string }

func (u *UnknownGroupError) Error() string {
	return fmt.Sprintf("group \"%s\" haven't found in regex", u.group)
}

// TypeNotParsableError returned when the type of struct field is not parsable
type TypeNotParsableError struct{ typ reflect.Type }

func (t *TypeNotParsableError) Error() string {
	return fmt.Sprintf("type \"%v\" is not parsable", t.typ)
}

// ParseError returned when the conversion to target struct field type has failed
type ParseError struct {
	group string
	err   error
}

func (p *ParseError) Error() string {
	return fmt.Sprintf("error parsing group \"%s\": %v", p.group, p.err)
}

// RequiredGroupIsEmpty returned when a required group is empty in the re match
type RequiredGroupIsEmpty struct {
	groupName string
	fieldName string
}

func (r *RequiredGroupIsEmpty) Error() string {
	return fmt.Sprintf("required regroup \"%s\" is empty for field \"%s\"", r.groupName, r.fieldName)
}
