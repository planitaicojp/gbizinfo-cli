package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cmdconfig "github.com/planitaicojp/gbizinfo-cli/cmd/config"
	"github.com/planitaicojp/gbizinfo-cli/cmd/update"
	cerrors "github.com/planitaicojp/gbizinfo-cli/internal/errors"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:           "gbizinfo",
	Short:         "gBizINFO REST API CLI",
	Long:          "gBizINFO（経済産業省 法人情報API）を操作するCLIツール",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringP("format", "f", "", "出力形式 (json/table/csv)")
	rootCmd.PersistentFlags().StringP("token", "t", "", "APIトークン")
	rootCmd.PersistentFlags().Bool("no-color", false, "色出力を無効化")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "詳細出力")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(cmdconfig.Cmd)
	rootCmd.AddCommand(update.Cmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(certificationCmd)
	rootCmd.AddCommand(commendationCmd)
	rootCmd.AddCommand(financeCmd)
	rootCmd.AddCommand(patentCmd)
	rootCmd.AddCommand(procurementCmd)
	rootCmd.AddCommand(subsidyCmd)
	rootCmd.AddCommand(workplaceCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cerrors.GetExitCode(err))
	}
}
