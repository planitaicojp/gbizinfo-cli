package update

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/model"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var Cmd = &cobra.Command{
	Use:   "update",
	Short: "期間指定で更新情報を取得",
	Long:  "指定した期間内に更新された法人情報を取得します。",
}

func init() {
	Cmd.AddCommand(hojinCmd)
	Cmd.AddCommand(certificationCmd)
	Cmd.AddCommand(commendationCmd)
	Cmd.AddCommand(financeCmd)
	Cmd.AddCommand(patentCmd)
	Cmd.AddCommand(procurementCmd)
	Cmd.AddCommand(subsidyCmd)
	Cmd.AddCommand(workplaceCmd)

	for _, cmd := range []*cobra.Command{
		hojinCmd, certificationCmd, commendationCmd, financeCmd,
		patentCmd, procurementCmd, subsidyCmd, workplaceCmd,
	} {
		addUpdateFlags(cmd)
	}
}

func updateParams(cmd *cobra.Command) model.UpdateParams {
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	page, _ := cmd.Flags().GetInt("page")
	return model.UpdateParams{From: from, To: to, Page: page}
}

func addUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().String("from", "", "開始日 (YYYY-MM-DD)")
	cmd.Flags().String("to", "", "終了日 (YYYY-MM-DD)")
	cmd.Flags().IntP("page", "p", 1, "ページ番号")
}

var hojinCmd = &cobra.Command{
	Use:     "hojin",
	Short:   "期間内に更新された法人情報を取得",
	Example: "  gbizinfo update hojin --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdateInfo(updateParams(cmd))
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

var certificationCmd = &cobra.Command{
	Use:     "certification",
	Short:   "期間内に更新された認定情報を取得",
	Example: "  gbizinfo update certification --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdateCertification(updateParams(cmd))
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

var commendationCmd = &cobra.Command{
	Use:     "commendation",
	Short:   "期間内に更新された表彰情報を取得",
	Example: "  gbizinfo update commendation --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdateCommendation(updateParams(cmd))
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

var financeCmd = &cobra.Command{
	Use:     "finance",
	Short:   "期間内に更新された財務情報を取得",
	Example: "  gbizinfo update finance --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdateFinance(updateParams(cmd))
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

var patentCmd = &cobra.Command{
	Use:     "patent",
	Short:   "期間内に更新された特許情報を取得",
	Example: "  gbizinfo update patent --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdatePatent(updateParams(cmd))
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

var procurementCmd = &cobra.Command{
	Use:     "procurement",
	Short:   "期間内に更新された調達情報を取得",
	Example: "  gbizinfo update procurement --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdateProcurement(updateParams(cmd))
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

var subsidyCmd = &cobra.Command{
	Use:     "subsidy",
	Short:   "期間内に更新された補助金情報を取得",
	Example: "  gbizinfo update subsidy --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdateSubsidy(updateParams(cmd))
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

var workplaceCmd = &cobra.Command{
	Use:     "workplace",
	Short:   "期間内に更新された職場情報を取得",
	Example: "  gbizinfo update workplace --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdateWorkplace(updateParams(cmd))
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
