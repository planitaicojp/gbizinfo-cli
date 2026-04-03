package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var certificationCmd = &cobra.Command{
	Use:     "certification [法人番号]",
	Short:   "届出・認定情報を取得",
	Long:    "指定した法人番号の届出・認定情報を取得します。",
	Example: "  gbizinfo certification 1234567890123",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		cn, err := cmdutil.CorporateNumberArg(cmd, args)
		if err != nil {
			return err
		}
		result, err := client.GetCertification(cn)
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	certificationCmd.Flags().StringP("corporate-number", "c", "", "法人番号")
}
