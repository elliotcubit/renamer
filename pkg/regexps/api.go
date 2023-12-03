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
	"regexp"
)

type Regexp[T any] struct {
	matcher  *regexp.Regexp
	defaults map[string]string
}

func MatchString[T any](expr, target string) (bool, error) {
	exVal := new(T)
	_, err := validateTarget(exVal)
	if err != nil {
		return false, &NotStructPtrError{}
	}
	pattern, err := Compile[T](expr)
	if err != nil {
		return false, err
	}
	return pattern.MatchString(target), nil
}

func CompileWithDefaults[T any](expr string, defaults map[string]string) (*Regexp[T], error) {
	// Verify the generic argument is acceptable
	exVal := new(T)
	targetRef, err := validateTarget(exVal)
	if err != nil {
		return nil, err
	}

	// Compile the pattern
	matcher, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	retv := &Regexp[T]{
		matcher:  matcher,
		defaults: defaults,
	}

	// Return the validated pattern
	return retv, retv.validateStruct(targetRef)
}

func Compile[T any](expr string) (*Regexp[T], error) {
	return CompileWithDefaults[T](expr, nil)
}

func MustCompileWithDefaults[T any](expr string, defaults map[string]string) *Regexp[T] {
	retv, err := CompileWithDefaults[T](expr, defaults)
	if err != nil {
		panic(err)
	}
	return retv
}

func MustCompile[T any](expr string) *Regexp[T] {
	retv, err := Compile[T](expr)
	if err != nil {
		panic(err)
	}
	return retv
}

func (r *Regexp[T]) MatchString(target string) bool {
	match := r.matcher.FindStringSubmatch(target)
	return len(match) > 0
}

func (r *Regexp[T]) Find(b []byte) *T {
	return r.FindString(string(b))
}

func (r *Regexp[T]) FindAll(b []byte) []*T {
	return r.FindAllString(string(b))
}

func (r *Regexp[T]) FindString(s string) *T {
	match := r.matcher.FindStringSubmatch(s)
	if match == nil {
		return nil
	}
	targetRef, _ := validateTarget(new(T))
	err := r.fillTarget(r.matchGroupMap(match), targetRef)
	if err != nil {
		return nil
	}
	retv := targetRef.Interface().(T)
	return &retv
}

func (r *Regexp[T]) FindAllString(s string) []*T {
	matches := r.matcher.FindAllStringSubmatch(s, -1)
	if matches == nil {
		return nil
	}
	retv := make([]*T, len(matches))
	for _, match := range matches {
		targetRef, _ := validateTarget(new(T))
		err := r.fillTarget(r.matchGroupMap(match), targetRef)
		if err != nil {
			return nil
		}
		nxt := targetRef.Interface().(T)
		retv = append(retv, &nxt)
	}
	return retv
}
