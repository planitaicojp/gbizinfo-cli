package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var subsidyCmd = &cobra.Command{
	Use:     "subsidy [法人番号]",
	Short:   "補助金情報を取得",
	Long:    "指定した法人番号の補助金情報を取得します。",
	Example: "  gbizinfo subsidy 1234567890123",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		cn, err := cmdutil.CorporateNumberArg(cmd, args)
		if err != nil {
			return err
		}
		result, err := client.GetSubsidy(cn)
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	subsidyCmd.Flags().StringP("corporate-number", "c", "", "法人番号")
}
