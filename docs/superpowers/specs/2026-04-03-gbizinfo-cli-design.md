# gbizinfo-cli 設計書

## 概要

gBizINFO（経済産業省 法人情報API）を操作するためのGo言語製CLIツール。法人基本情報・財務情報・補助金受給履歴等を検索・取得する。

- **対象API**: gBizINFO REST API (`https://info.gbiz.go.jp/hojin`)
- **認証**: `X-hojinInfo-api-token` ヘッダー（アカウント登録で無料取得）
- **Swagger Spec**: `https://info.gbiz.go.jp/hojin/v2/api-docs`（Swagger 2.0 JSON）

## 設計方針

- Go言語で実装（参照: conoha-cli の構造パターン）
- 動詞中心のサブコマンド構成
- 出力形式: JSON（デフォルト）、Table、CSV
- コマンド・フラグ名は英語、ヘルプメッセージは日本語
- AIエージェント（Claude Code等）からの利用を想定
- Rate Limit（429）は自動リトライせず、エラーメッセージのみ表示

## 対象API エンドポイント（全17個）

### 法人情報取得（9個）

| エンドポイント | 説明 |
|---------------|------|
| `GET /v1/hojin` | 法人検索 |
| `GET /v1/hojin/{corporate_number}` | 法人基本情報 |
| `GET /v1/hojin/{corporate_number}/certification` | 届出・認定情報 |
| `GET /v1/hojin/{corporate_number}/commendation` | 表彰情報 |
| `GET /v1/hojin/{corporate_number}/finance` | 財務情報 |
| `GET /v1/hojin/{corporate_number}/patent` | 特許情報 |
| `GET /v1/hojin/{corporate_number}/procurement` | 調達情報 |
| `GET /v1/hojin/{corporate_number}/subsidy` | 補助金情報 |
| `GET /v1/hojin/{corporate_number}/workplace` | 職場情報 |

### 期間指定更新情報（8個）

| エンドポイント | 説明 |
|---------------|------|
| `GET /v1/hojin/updateInfo` | 期間内更新法人情報 |
| `GET /v1/hojin/updateInfo/certification` | 期間内更新認定情報 |
| `GET /v1/hojin/updateInfo/commendation` | 期間内更新表彰情報 |
| `GET /v1/hojin/updateInfo/finance` | 期間内更新財務情報 |
| `GET /v1/hojin/updateInfo/patent` | 期間内更新特許情報 |
| `GET /v1/hojin/updateInfo/procurement` | 期間内更新調達情報 |
| `GET /v1/hojin/updateInfo/subsidy` | 期間内更新補助金情報 |
| `GET /v1/hojin/updateInfo/workplace` | 期間内更新職場情報 |

## サブコマンド体系

```
gbizinfo
├── search          # GET /v1/hojin — 法人検索
├── get             # GET /v1/hojin/{corp} — 基本情報
├── certification   # GET /v1/hojin/{corp}/certification
├── commendation    # GET /v1/hojin/{corp}/commendation
├── finance         # GET /v1/hojin/{corp}/finance
├── patent          # GET /v1/hojin/{corp}/patent
├── procurement     # GET /v1/hojin/{corp}/procurement
├── subsidy         # GET /v1/hojin/{corp}/subsidy
├── workplace       # GET /v1/hojin/{corp}/workplace
├── update          # updateInfo 系グループコマンド
│   ├── hojin           # GET /v1/hojin/updateInfo
│   ├── certification
│   ├── commendation
│   ├── finance
│   ├── patent
│   ├── procurement
│   ├── subsidy
│   └── workplace
├── config          # 設定管理
│   ├── init            # 初期設定（インタラクティブ）
│   ├── set             # 個別値設定
│   └── show            # 現在設定表示（トークンはマスク）
├── version         # バージョン表示
└── completion      # シェル自動補完
```

### 法人番号指定コマンド（get, finance, subsidy 等）

- 第1引数で法人番号を受取: `gbizinfo finance 1234567890123`
- `--corporate-number` / `-c` フラグでも指定可能

### search コマンドフラグ

| フラグ | 短縮 | 説明 |
|--------|------|------|
| `--name` | `-n` | 法人名 |
| `--address` | | 所在地 |
| `--corporate-number` | `-c` | 法人番号 |
| `--page` | `-p` | ページ番号 |
| `--limit` | `-l` | 表示件数 |

### update 系共通フラグ

| フラグ | 説明 |
|--------|------|
| `--from` | 開始日（YYYY-MM-DD） |
| `--to` | 終了日（YYYY-MM-DD） |
| `--page` / `-p` | ページ番号 |

### グローバルフラグ（root）

| フラグ | 短縮 | 説明 | デフォルト |
|--------|------|------|-----------|
| `--format` | `-f` | 出力形式（json/table/csv） | json |
| `--token` | `-t` | APIトークン（環境変数オーバーライド） | — |
| `--no-color` | | 色無効化 | false |
| `--verbose` | `-v` | 詳細出力（リクエスト/レスポンスヘッダー） | false |

## APIクライアント

### 構造

```go
// internal/api/client.go
type Client struct {
    HTTP    *http.Client
    BaseURL string   // https://info.gbiz.go.jp/hojin
    Token   string   // X-hojinInfo-api-token
}
```

- 単一クライアントに全メソッドを配置（gBizINFOは単一サービス）
- ファイル分離: `client.go`（ベース）+ `hojin.go`（法人取得9個）+ `update.go`（期間指定8個）

### ヘッダー

- `X-hojinInfo-api-token`: 認証トークン
- `Accept: application/json`
- `User-Agent: gbizinfo-cli/{version}`

### エラー処理

- HTTP >= 400 で構造化エラーをパース
- 401 → `AuthError`
- 404 → `NotFoundError`
- 429 → `RateLimitError`（メッセージのみ、自動リトライなし）
- その他 → `APIError`

### デバッグ

- `--verbose` 時にリクエスト/レスポンスヘッダーをstderrに出力

### メソッドパターン

```go
func (c *Client) Search(params SearchParams) (*HojinResponse, error)
func (c *Client) Get(corporateNumber string) (*Hojin, error)
func (c *Client) GetFinance(corporateNumber string) (*FinanceResponse, error)
func (c *Client) GetCertification(corporateNumber string) (*CertificationResponse, error)
func (c *Client) GetCommendation(corporateNumber string) (*CommendationResponse, error)
func (c *Client) GetPatent(corporateNumber string) (*PatentResponse, error)
func (c *Client) GetProcurement(corporateNumber string) (*ProcurementResponse, error)
func (c *Client) GetSubsidy(corporateNumber string) (*SubsidyResponse, error)
func (c *Client) GetWorkplace(corporateNumber string) (*WorkplaceResponse, error)
func (c *Client) GetUpdateInfo(params UpdateParams) (*UpdateResponse, error)
func (c *Client) GetUpdateCertification(params UpdateParams) (*UpdateResponse, error)
// ... 他のUpdate系メソッドも同様
```

## 設定管理

### ファイル

- 場所: `~/.config/gbizinfo/config.yaml`（`XDG_CONFIG_HOME` 尊重）
- パーミッション: `0600`（トークン保護）

### 設定ファイル形式

```yaml
token: "your-api-token-here"
defaults:
  format: json
```

### 優先順位（高い順）

1. CLIフラグ（`--token`, `--format`）
2. 環境変数（`GBIZINFO_TOKEN`, `GBIZINFO_FORMAT`）
3. 設定ファイル
4. デフォルト値（format: json）

### config サブコマンド

- `gbizinfo config init` — トークン入力プロンプト、ファイル生成
- `gbizinfo config set token <value>` — 個別値設定
- `gbizinfo config show` — 現在設定表示（トークンはマスク表示）

## データモデル

### ファイル構成

```
internal/model/
├── hojin.go         # Hojin（基本情報）、HojinResponse（検索結果+ページネーション）
├── certification.go # Certification
├── commendation.go  # Commendation
├── finance.go       # Finance
├── patent.go        # Patent
├── procurement.go   # Procurement
├── subsidy.go       # Subsidy
├── workplace.go     # Workplace
└── update.go        # UpdateResponse（期間指定取得共通）
```

### ページネーション

```go
type PageInfo struct {
    TotalCount int `json:"totalCount"`
    TotalPage  int `json:"totalPage"`
    PageNumber int `json:"pageNumber"`
}

type HojinResponse struct {
    PageInfo
    Corporations []Hojin `json:"hojin-infos"`
}
```

- JSONタグでAPIレスポンスをマッピング
- 同じタグをTable/CSVカラム名としても使用

## 出力フォーマット

### インターフェース

```go
// internal/output/formatter.go
type Formatter interface {
    Format(w io.Writer, data any) error
}

func New(format string) Formatter  // ファクトリ関数
```

### 実装

| フォーマッタ | ファイル | 説明 |
|------------|---------|------|
| JSON | `json.go` | 2スペースインデント、`json.Encoder` |
| Table | `table.go` | `text/tabwriter`、リフレクションでフィールド抽出 |
| CSV | `csv.go` | `encoding/csv`、ヘッダー行+データ行 |

### 使用パターン

```go
format := cmdutil.GetFormat(cmd)
formatter := output.New(format)
formatter.Format(os.Stdout, data)
```

## エラー処理と終了コード

### カスタムエラー型

```go
AuthError        // 401 — トークン未設定・無効
NotFoundError    // 404 — 法人番号なし
RateLimitError   // 429 — API制限超過
APIError         // その他 4xx/5xx
ValidationError  // フラグ・引数検証失敗
```

### 終了コード

| コード | 意味 |
|-------|------|
| 0 | 成功 |
| 1 | 一般エラー |
| 2 | 認証失敗 |
| 3 | Not Found |
| 4 | 入力検証エラー |
| 5 | APIエラー |

- `ExitCoder` インターフェースパターン（conoha-cli準拠）
- `root.go` の `Execute()` でエラーをキャッチして終了コード決定
- エラーメッセージは日本語でstderrに出力

## ディレクトリ構造

```
gbizinfo-cli/
├── main.go
├── go.mod
├── Makefile
├── .golangci.yml
├── README.md
├── cmd/
│   ├── root.go
│   ├── version.go
│   ├── completion.go
│   ├── search.go
│   ├── get.go
│   ├── certification.go
│   ├── commendation.go
│   ├── finance.go
│   ├── patent.go
│   ├── procurement.go
│   ├── subsidy.go
│   ├── workplace.go
│   ├── update/
│   │   └── update.go        # グループ + 8サブコマンド
│   ├── config/
│   │   └── config.go        # init, set, show
│   └── cmdutil/
│       └── cmdutil.go        # クライアント生成、フォーマット取得
├── internal/
│   ├── api/
│   │   ├── client.go         # ベースHTTPクライアント
│   │   ├── hojin.go          # 法人取得メソッド（9個）
│   │   └── update.go         # 期間指定メソッド（8個）
│   ├── config/
│   │   └── config.go         # YAML設定 + 環境変数
│   ├── model/
│   │   ├── hojin.go
│   │   ├── certification.go
│   │   ├── commendation.go
│   │   ├── finance.go
│   │   ├── patent.go
│   │   ├── procurement.go
│   │   ├── subsidy.go
│   │   ├── workplace.go
│   │   └── update.go
│   ├── output/
│   │   ├── formatter.go
│   │   ├── json.go
│   │   ├── table.go
│   │   └── csv.go
│   └── errors/
│       └── errors.go
└── test/
    └── fixtures/             # テスト用JSONレスポンス
```

## 技術スタック

| 項目 | 選定 |
|------|------|
| 言語 | Go |
| CLIフレームワーク | github.com/spf13/cobra |
| YAML | gopkg.in/yaml.v3 |
| テーブル出力 | text/tabwriter（標準ライブラリ） |
| CSV出力 | encoding/csv（標準ライブラリ） |
| HTTP | net/http（標準ライブラリ） |
| JSON | encoding/json（標準ライブラリ） |
| テスト | net/http/httptest（標準ライブラリ） |
| リンター | golangci-lint |
