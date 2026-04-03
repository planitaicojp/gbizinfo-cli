# gbizinfo-cli

gBizINFO（経済産業省 法人情報API）を操作するCLIツール。

## インストール

### Homebrew (macOS / Linux)

```bash
brew install planitaicojp/tap/gbizinfo-cli
```

### Scoop (Windows)

```powershell
scoop bucket add planitaicojp https://github.com/planitaicojp/scoop-bucket.git
scoop install gbizinfo-cli
```

### Go

```bash
go install github.com/planitaicojp/gbizinfo-cli@latest
```

### リリースバイナリ

[GitHub Releases](https://github.com/planitaicojp/gbizinfo-cli/releases) からダウンロード:

```bash
# Linux (amd64)
curl -sL https://github.com/planitaicojp/gbizinfo-cli/releases/latest/download/gbizinfo-cli_Linux_amd64.tar.gz | tar xz
sudo mv gbizinfo /usr/local/bin/

# macOS (Apple Silicon)
curl -sL https://github.com/planitaicojp/gbizinfo-cli/releases/latest/download/gbizinfo-cli_Darwin_arm64.tar.gz | tar xz
sudo mv gbizinfo /usr/local/bin/

# Windows (amd64) — PowerShell
Invoke-WebRequest -Uri https://github.com/planitaicojp/gbizinfo-cli/releases/latest/download/gbizinfo-cli_Windows_amd64.zip -OutFile gbizinfo.zip
Expand-Archive gbizinfo.zip -DestinationPath .
```

## 初期設定

[gBizINFO](https://info.gbiz.go.jp/) でアカウント登録後、APIトークンを取得してください。

```bash
gbizinfo config init
```

環境変数でも設定できます:

```bash
export GBIZINFO_TOKEN=your-token-here
```

## 使い方

### 法人検索

```bash
gbizinfo search --name トヨタ
gbizinfo search --name トヨタ --address 愛知県
gbizinfo search -c 1234567890123
```

### 法人情報取得

```bash
gbizinfo get 1234567890123          # 基本情報
gbizinfo finance 1234567890123      # 財務情報
gbizinfo subsidy 1234567890123      # 補助金情報
gbizinfo patent 1234567890123       # 特許情報
gbizinfo procurement 1234567890123  # 調達情報
gbizinfo certification 1234567890123 # 届出・認定情報
gbizinfo commendation 1234567890123 # 表彰情報
gbizinfo workplace 1234567890123    # 職場情報
```

### 期間指定更新情報

```bash
gbizinfo update hojin --from 2024-01-01 --to 2024-01-31
gbizinfo update finance --from 2024-01-01 --to 2024-12-31
```

### 出力形式

```bash
gbizinfo search --name トヨタ -f table
gbizinfo search --name トヨタ -f csv
gbizinfo search --name トヨタ -f json   # デフォルト
```

## 設定

設定ファイル: `~/.config/gbizinfo/config.yaml`

```yaml
token: "your-api-token"
defaults:
  format: json
```

### 環境変数

| 変数名 | 説明 |
|--------|------|
| `GBIZINFO_TOKEN` | APIトークン |
| `GBIZINFO_FORMAT` | デフォルト出力形式 (json/table/csv) |
| `GBIZINFO_CONFIG_DIR` | 設定ディレクトリのパス |

### 優先順位

CLIフラグ > 環境変数 > 設定ファイル > デフォルト値

## 対象API

[gBizINFO REST API](https://info.gbiz.go.jp/api/index.html) — 経済産業省が提供する法人情報API

## ライセンス

MIT
