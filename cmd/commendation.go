package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var commendationCmd = &cobra.Command{
	Use:     "commendation [法人番号]",
	Short:   "表彰情報を取得",
	Long:    "指定した法人番号の表彰情報を取得します。",
	Example: "  gbizinfo commendation 1234567890123",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		cn, err := cmdutil.CorporateNumberArg(cmd, args)
		if err != nil {
			return err
		}
		result, err := client.GetCommendation(cn)
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	commendationCmd.Flags().StringP("corporate-number", "c", "", "法人番号")
}
