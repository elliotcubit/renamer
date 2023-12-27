package file

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/elliotcubit/renamer/pkg/regexps"
)

var errCantInfer = errors.New("cannot infer pattern of file names")

var rawPatterns = []string{
	`(?P<name>[^-]+) - \[(?P<season>\d+)x(?P<episode>\d+)\] - (?P<title>.*)\....`,
	`(?P<name>[^\.]+)\.S(?P<season>\d+)E(?P<episode>\d\d)\.(?P<title>.*)\.?\d+p.*\....`,
}

var patterns []*regexps.Regexp[Match]

func init() {
	patterns = make([]*regexps.Regexp[Match], len(rawPatterns))
	for i, v := range rawPatterns {
		patterns[i] = regexps.MustCompile[Match](v)
	}
}

// Check if all files in a directory match a particular pattern.
// If they do, we can reasonable say that's what the user wants.
func InferPattern(
	fsys fs.FS,
	dir string,
) (*regexps.Regexp[Match], error) {

	matched := make([]bool, len(patterns))
	for i := range matched {
		matched[i] = true
	}

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		_, file := filepath.Split(path)
		for i, v := range patterns {
			matched[i] = matched[i] && (v.MatchString(file) || isOkay(file))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for i, v := range matched {
		if v {
			fmt.Printf("Using detected pattern %q\n", rawPatterns[i])
			return patterns[i], nil
		}
	}

	return nil, errCantInfer
}

// It is okay if these files don't match the patterns :)
var blocklist = []string{
	"torrent",
	"read me",
	"readme",
}

func isOkay(fname string) bool {
	f := strings.ToLower(fname)
	for _, v := range blocklist {
		if strings.Contains(f, v) {
			return true
		}
	}
	return false
}
