package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var workplaceCmd = &cobra.Command{
	Use:     "workplace [法人番号]",
	Short:   "職場情報を取得",
	Long:    "指定した法人番号の職場情報（従業員数、女性比率等）を取得します。",
	Example: "  gbizinfo workplace 1234567890123",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		cn, err := cmdutil.CorporateNumberArg(cmd, args)
		if err != nil {
			return err
		}
		result, err := client.GetWorkplace(cn)
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	workplaceCmd.Flags().StringP("corporate-number", "c", "", "法人番号")
}
