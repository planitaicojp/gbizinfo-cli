package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	iconfig "github.com/planitaicojp/gbizinfo-cli/internal/config"
)

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "設定を管理",
}

func init() {
	Cmd.AddCommand(initCmd)
	Cmd.AddCommand(setCmd)
	Cmd.AddCommand(showCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初期設定を行う",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprint(os.Stderr, "gBizINFO APIトークンを入力してください: ")
		reader := bufio.NewReader(os.Stdin)
		token, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("入力の読み取りに失敗: %w", err)
		}
		token = strings.TrimSpace(token)
		if token == "" {
			return fmt.Errorf("トークンが入力されませんでした")
		}

		cfg, err := iconfig.Load()
		if err != nil {
			return err
		}
		cfg.Token = token
		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "設定を保存しました: %s/config.yaml\n", iconfig.DefaultConfigDir())
		return nil
	},
}

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "設定値を変更",
	Long:  "設定値を変更します。キー: token, format",
	Args:  cmdutil.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]

		cfg, err := iconfig.Load()
		if err != nil {
			return err
		}

		switch key {
		case "token":
			cfg.Token = value
		case "format":
			cfg.Defaults.Format = value
		default:
			return fmt.Errorf("不明な設定キー: %s (使用可能: token, format)", key)
		}

		return cfg.Save()
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "現在の設定を表示",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := iconfig.Load()
		if err != nil {
			return err
		}

		fmt.Printf("設定ディレクトリ: %s\n", iconfig.DefaultConfigDir())
		fmt.Printf("トークン:         %s\n", iconfig.MaskToken(cfg.Token))
		fmt.Printf("出力形式:         %s\n", cfg.Defaults.Format)
		return nil
	},
}
