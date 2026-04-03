package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DefaultFormat = "json"
	configFile    = "config.yaml"
)

const (
	EnvConfigDir = "GBIZINFO_CONFIG_DIR"
	EnvToken     = "GBIZINFO_TOKEN"
	EnvFormat    = "GBIZINFO_FORMAT"
)

type Config struct {
	Token    string   `yaml:"token"`
	Defaults Defaults `yaml:"defaults"`
}

type Defaults struct {
	Format string `yaml:"format"`
}

func DefaultConfigDir() string {
	if d := os.Getenv(EnvConfigDir); d != "" {
		return d
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "gbizinfo")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "gbizinfo")
}

func Load() (*Config, error) {
	dir := DefaultConfigDir()
	path := filepath.Join(dir, configFile)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultConfig(), nil
		}
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("設定ファイルの解析に失敗: %w", err)
	}
	if cfg.Defaults.Format == "" {
		cfg.Defaults.Format = DefaultFormat
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	dir := DefaultConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("設定ディレクトリの作成に失敗: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("設定のシリアライズに失敗: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, configFile), data, 0600)
}

func EnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func MaskToken(token string) string {
	if token == "" {
		return ""
	}
	if len(token) <= 4 {
		return strings.Repeat("*", len(token))
	}
	return token[:4] + strings.Repeat("*", 6)
}

func defaultConfig() *Config {
	return &Config{
		Defaults: Defaults{Format: DefaultFormat},
	}
}
