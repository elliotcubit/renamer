package cmd

import (
	"fmt"
	"os"

	"github.com/elliotcubit/renamer/pkg/file"
	"github.com/elliotcubit/renamer/pkg/regexps"
	"github.com/spf13/cobra"
)

const (
	PatternFlagName  = "pattern"
	TemplateFlagName = "template"
	DirFlagName      = "dir"
	DryRunFlagName   = "dry-run"
	OutputFlagName   = "output-template"
)

func init() {
	rootCmd.PersistentFlags().StringP(PatternFlagName, "p", "", "Pattern of files to pick up")
	rootCmd.PersistentFlags().StringP(DirFlagName, "d", "", "Directory to check")
	rootCmd.PersistentFlags().Bool(DryRunFlagName, false, "Do not modify any files; instead, print what would be done")
	rootCmd.PersistentFlags().StringP(OutputFlagName, "o", "{{ .ShowName }} s{{ .Season }}e{{ .Episode }} - {{ .Title }}", "The template to rename files to, not including any file extension")
	cobra.MarkFlagRequired(rootCmd.PersistentFlags(), PatternFlagName)

	for k, v := range defaultArgs {
		rootCmd.PersistentFlags().String(k, "", v)
	}
}

var defaultArgs = map[string]string{
	"name":   "The name of the show",
	"season": "The season the episode is in",
}

var rootCmd = &cobra.Command{
	Use:   "renamer",
	Short: "renamer renames files to a standard format.",
	Run: func(cmd *cobra.Command, args []string) {
		defaults := make(map[string]string, 0)
		for k := range defaultArgs {
			if v := cmd.Flag(k).Value.String(); v != "" {
				defaults[k] = v
			}
		}
		pattern, err := regexps.CompileWithDefaults[file.Match](cmd.Flag(PatternFlagName).Value.String(), defaults)
		if err != nil {
			fmt.Printf("bad pattern: %v\n", err)
			os.Exit(1)
		}

		dir := cmd.Flag(DirFlagName).Value.String()

		fs := os.DirFS(dir)
		err = file.RenameAllFiles(
			fs,
			dir,
			pattern,
			cmd.Flag(DryRunFlagName).Changed,
			cmd.Flag(OutputFlagName).Value.String(),
		)
		if err != nil {
			fmt.Printf("rename: %v", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
