package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "シェル補完スクリプトを生成",
	Long: `シェル補完スクリプトを生成します。

使用例:
  # Bash
  gbizinfo completion bash > /etc/bash_completion.d/gbizinfo

  # Zsh
  gbizinfo completion zsh > "${fpath[1]}/_gbizinfo"

  # Fish
  gbizinfo completion fish > ~/.config/fish/completions/gbizinfo.fish`,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cmdutil.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return nil
		}
	},
}
