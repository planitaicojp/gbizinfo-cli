package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/model"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "法人を検索",
	Long:  "gBizINFOに登録された法人を検索します。",
	Example: `  gbizinfo search --name トヨタ
  gbizinfo search --name トヨタ --address 愛知県
  gbizinfo search -c 1234567890123`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		address, _ := cmd.Flags().GetString("address")
		cn, _ := cmd.Flags().GetString("corporate-number")
		page, _ := cmd.Flags().GetInt("page")
		limit, _ := cmd.Flags().GetInt("limit")

		result, err := client.Search(model.SearchParams{
			Name:            name,
			Address:         address,
			CorporateNumber: cn,
			Page:            page,
			Limit:           limit,
		})
		if err != nil {
			return err
		}

		format := cmdutil.GetFormat(cmd)
		if format == "json" {
			return output.New(format).Format(os.Stdout, result)
		}
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	searchCmd.Flags().StringP("name", "n", "", "法人名")
	searchCmd.Flags().String("address", "", "所在地")
	searchCmd.Flags().StringP("corporate-number", "c", "", "法人番号")
	searchCmd.Flags().IntP("page", "p", 1, "ページ番号")
	searchCmd.Flags().IntP("limit", "l", 0, "表示件数")
}
