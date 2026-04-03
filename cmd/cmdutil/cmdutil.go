package cmdutil

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/internal/api"
	"github.com/planitaicojp/gbizinfo-cli/internal/config"
	cerrors "github.com/planitaicojp/gbizinfo-cli/internal/errors"
)

const defaultBaseURL = "https://info.gbiz.go.jp/hojin"

func NewClient(cmd *cobra.Command) (*api.Client, error) {
	token, _ := cmd.Flags().GetString("token")
	if token == "" {
		token = config.EnvOr(config.EnvToken, "")
	}
	if token == "" {
		cfg, err := config.Load()
		if err != nil {
			return nil, err
		}
		token = cfg.Token
	}
	if token == "" {
		return nil, &cerrors.AuthError{Message: "APIトークンが設定されていません。gbizinfo config init を実行してください"}
	}

	verbose, _ := cmd.Flags().GetBool("verbose")
	client := api.NewClient(defaultBaseURL, token)
	client.Verbose = verbose
	return client, nil
}

func GetFormat(cmd *cobra.Command) string {
	format, _ := cmd.Flags().GetString("format")
	if format != "" {
		return format
	}
	if f := config.EnvOr(config.EnvFormat, ""); f != "" {
		return f
	}
	cfg, err := config.Load()
	if err != nil {
		return config.DefaultFormat
	}
	if cfg.Defaults.Format != "" {
		return cfg.Defaults.Format
	}
	return config.DefaultFormat
}

func ExactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("%d個の引数が必要です（%d個指定されました）\n\n使い方:\n  %s", n, len(args), cmd.UseLine())
		}
		return nil
	}
}

func CorporateNumberArg(cmd *cobra.Command, args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}
	cn, _ := cmd.Flags().GetString("corporate-number")
	if cn != "" {
		return cn, nil
	}
	return "", &cerrors.ValidationError{Field: "corporate_number", Message: "法人番号を指定してください"}
}
