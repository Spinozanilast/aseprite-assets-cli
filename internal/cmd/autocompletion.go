package cmd

import (
	"sync"

	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/utils/files"
)

type AutocompletionFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

func GenerateFilesAutoCompletions(folders []string, extensions []string) AutocompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		wg := new(sync.WaitGroup)
		var completions []string

		for _, dir := range folders {
			wg.Add(1)
			go func(folder string) {
				defer wg.Done()

				fs, err := files.FindFilesOfExtensionsRecursiveFlatten(dir, extensions...)
				if err == nil {
					completions = append(completions, fs...)
				}
			}(dir)
		}
		wg.Wait()
		return completions, cobra.ShellCompDirectiveNoFileComp
	}
}
