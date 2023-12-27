package file

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

func wrapNamesInFS(fnames []string) fs.FS {
	retv := make(map[string]*fstest.MapFile, len(fnames))
	for _, v := range fnames {
		retv[v] = &fstest.MapFile{}
	}
	return fstest.MapFS(retv)
}

func TestInferFilenames(t *testing.T) {
	type testcase struct {
		name    string
		matches bool
	}

	tests := []testcase{
		{"House - [4x04] - Guardian Angels.mp4", true},
		{"You.S02E05.Dont.720p.mkv", true},
		{"You.S02E05.Dont1024p.skv", true},
	}

	for _, test := range tests {
		fs := wrapNamesInFS([]string{test.name})
		pat, err := InferPattern(fs, ".")
		if test.matches {
			if err != nil {
				t.Fatalf("infer on %q: %v", test.name, err)
			}
			if pat == nil {
				t.Fatalf("infer: nil pattern for %q", test.name)
			}
		} else {
			if pat != nil {
				t.Fatalf("infer: non-nil pattern for %q", test.name)
			}
			if err == nil {
				t.Fatalf("infer: nil err for %q", test.name)
			}
		}

	}
}
