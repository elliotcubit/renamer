# Renamer

`renamer`` renames files. It is intended to be used for renamed content
to standard formats recognized by Plex.

## Usage

```
$ go run . -h
renamer renames files to a standard format.

Usage:
  renamer [flags]

Flags:
  -d, --dir string               Directory to check (default ".")
      --dry-run                  Do not modify any files; instead, print what would be done
  -h, --help                     help for renamer
      --name string              The name of the show
  -o, --output-template string   The template to rename files to, not including any file extension (default "{{ .ShowName }} s{{ .Season }}e{{ .Episode }} - {{ .Title }}")
  -p, --pattern string           Pattern of files to pick up
      --season string            The season the episode is in
```

The `--pattern` is a regular expression using named capture groups with the keys `episode`, `season`, `name` and `title`.

If _all_ files in the target directory match one of the included patterns (and _the same_ pattern), that
pattern can be used without providing the `--pattern` argument. There are a few detectable file patterns,
which will hopefully be expanded later.

The `--output-pattern` is a go template using those variables.

The `name` and `season` can be fixed by arugments, in which case they are not required in the input `--pattern`.

## License

### My original work

MIT

### Open-source license information

The `regexps` package is licensed under Apache 2.0, being adapted from [regroup](https://github.com/oriser/regroup), adding the ability to provide default arguments and an "exists" struct tag, which differs from the "requires" struct tag in that it requires the capture group to _exist_, but does not require it to be _populated_.

This package also allows creating a `Regexp` opject with a generic argument, instead of passing a pointer to a struct, and changes the public-facing API to be more in-line with the stdlib `regexp` package.

I have also added the requisite copyright notices to the package, which were not present in the original distribution.