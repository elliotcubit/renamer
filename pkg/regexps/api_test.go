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

import "testing"

func TestMatching(t *testing.T) {
	type s struct {
		Foo string `regexps:"foo"`
	}

	ok, err := MatchString[s]("(?P<foo>.*)", "hello")
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("did not match")
	}
}

func TestFinding(t *testing.T) {
	type s struct {
		Foo string `regexps:"foo,required"`
	}

	pattern := MustCompile[s]("(?P<foo>.*)")

	match := pattern.FindString("hello")
	if match == nil {
		t.Error("no match")
	} else {
		if match.Foo != "hello" {
			t.Errorf("wrong value, got: %q", match.Foo)
		}
	}
}

func TestDefaults(t *testing.T) {
	type s struct {
		Foo string `regexps:"foo,required"`
	}

	pattern := MustCompileWithDefaults[s](".*", map[string]string{"foo": "hello"})

	match := pattern.FindString("asdf")
	if match == nil {
		t.Error("no match")
	} else {
		if match.Foo != "hello" {
			t.Errorf("wrong value, got: %q", match.Foo)
		}
	}
}

func TestCompile(t *testing.T) {
	type s struct {
		Foo string `regexps:"foo"`
		Bar string `regexps:"bar,required"`
		Baz struct {
			Foo string `regexps:"foo2,required"`
		}
	}

	type testcase struct {
		pattern  string
		defaults map[string]string
		success  bool
	}

	tests := []testcase{
		{
			"(?P<foo>.*)",
			nil,
			false,
		},
		{
			"(?P<foo>.*)(?P<bar>.*)",
			nil,
			false,
		},
		{
			"(?P<foo>.*)(?P<bar>.*)(?P<foo2>.*)",
			nil,
			true,
		},
		{
			"(?P<foo>.*)(?P<bar>.*)",
			map[string]string{"foo2": "bar"},
			true,
		},
		{
			"(?P<foo>.*)?(?P<bar>.*)?(?P<foo2>.*)?",
			nil,
			true,
		},
	}

	for _, test := range tests {
		_, err := CompileWithDefaults[s](test.pattern, test.defaults)
		if (err == nil) != test.success {
			t.Errorf("fail on %q: err=%v", test.pattern, err)
		}
	}
}
