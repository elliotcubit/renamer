package file

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/elliotcubit/renamer/pkg/regexps"
)

type Match struct {
	ShowName string `regexps:"name,required"`
	Season   int    `regexps:"season,required"`
	Episode  int    `regexps:"episode,required"`
	Title    string `regexps:"title,required"`
}

func RenameAllFiles(
	fsys fs.FS,
	dir string,
	pattern *regexps.Regexp[Match],
	dry bool,
	outputTemplate string,
) error {
	if dry {
		fmt.Printf("In %q, would:\n", dir)
	}
	renamedSomething := false

	tmpl, err := template.New("output").Parse(outputTemplate)
	if err != nil {
		return fmt.Errorf("bad template: %w", err)
	}

	buf := new(strings.Builder)

	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		defer buf.Reset()
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		dir2, file := filepath.Split(path)
		ext := filepath.Ext(file)
		fullPath := filepath.Join(dir, path)

		match := pattern.FindString(file)
		if match != nil {
			err := tmpl.Execute(buf, match)
			if err != nil {
				return fmt.Errorf("apply template: %w", err)
			}

			newFile := buf.String()
			newFile += ext

			newPath := filepath.Join(dir, dir2, newFile)

			if dry {
				fmt.Printf("  Rename %q -> %q\n", path, newFile)
				renamedSomething = true
			} else {
				err := os.Rename(fullPath, newPath)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if !renamedSomething {
		fmt.Printf("  Do nothing (no files matched?)\n")
	}

	return err
}
