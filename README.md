# envault

.env ファイルを暗号化し、環境変数として安全に管理するためのCLIツール

## 概要

envaultは、.envファイルを暗号化して.env.vaultedファイルとして保存し、必要なときに復号化して環境変数をエクスポートするためのツールです。セキュリティを向上させながら、環境変数の管理を簡単に行うことができます。

## 機能

- .envファイルの暗号化
- 暗号化されたファイルからの環境変数のエクスポート
- エクスポートした環境変数のアンセット

## インストール

```bash
# リポジトリをクローン
git clone https://github.com/uzulla/envault.git
cd envault

# 依存関係を整理
go mod tidy

# ビルド
go build -o envault cmd/envault/main.go
```

または、Go Modulesを使用してインストール：

```bash
go install github.com/uzulla/envault/cmd/envault@latest
```

## 使用方法

### 暗号化

```bash
# .envファイルを暗号化して.env.vaultedファイルを作成
envault .env

# 出力ファイル名を指定して暗号化
envault .env --file /path/to/output.vaulted
```

### 環境変数のエクスポート

```bash
# .env.vaultedファイルから環境変数をエクスポート（重要：evalまたはsourceで実行する必要があります）
eval $(envault export)
# または
source <(envault export)

# 特定の暗号化ファイルから環境変数をエクスポート
eval $(envault export --file /path/to/custom.vaulted)

# stdinからパスワードを読み込んでエクスポート
echo "password" | envault export --password-stdin | source
```

**注意**: `envault export`コマンドはシェルスクリプトを出力するだけで、環境変数を直接設定しません。環境変数を実際に設定するには、上記のように`eval`または`source`コマンドを使用する必要があります。

### 環境変数のアンセット

```bash
# .env.vaultedファイルに記載された環境変数をアンセット（重要：evalまたはsourceで実行する必要があります）
eval $(envault unset)
# または
source <(envault unset)

# 特定の暗号化ファイルに記載された環境変数をアンセット
eval $(envault unset --file /path/to/custom.vaulted)

# stdinからパスワードを読み込んでアンセット
echo "password" | envault unset --password-stdin | source
```

**注意**: `envault unset`コマンドもシェルスクリプトを出力するだけで、環境変数を直接アンセットしません。環境変数を実際にアンセットするには、上記のように`eval`または`source`コマンドを使用する必要があります。

### ヘルプとバージョン情報

```bash
# ヘルプを表示
envault help

# バージョン情報を表示
envault version
```

## セキュリティ

- AES-256-GCMによる強力な暗号化
- Argon2idによる安全なパスワード派生関数
- 暗号化されたファイルからは環境変数のキー名や値を推測できない

## 詳細なドキュメント

- [設計ドキュメント](./Docs/design.md)
- [実装計画](./Docs/implementation_plan.md)
- [テスト手順](./QA/README.md)

## 動作環境

- Linux または macOS (bash)

## ライセンス

MIT
