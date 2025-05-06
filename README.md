# envault

.env ファイルを暗号化し、環境変数として安全に管理するためのCLIツール

## 概要

envaultは、.envファイルを暗号化して.env.vaultedファイルとして保存し、必要なときに復号化して環境変数をエクスポートするためのツールです。セキュリティを向上させながら、環境変数の管理を簡単に行うことができます。

## 機能

- .envファイルの暗号化
- 暗号化されたファイルからの環境変数のエクスポート
- エクスポートした環境変数のアンセット
- TUIを使用した環境変数の選択的なエクスポート/アンセット

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
envault .env -f /path/to/output.vaulted
# または
envault .env --file /path/to/output.vaulted
```

### 環境変数のエクスポート

#### 従来の方法（シェルスクリプト評価）

```bash
# .env.vaultedファイルから環境変数をエクスポート（重要：evalまたはsourceで実行する必要があります）
eval $(envault export -o)
# または
eval $(envault export --output-script-only)
# または
source <(envault export --output-script-only)

# 特定の暗号化ファイルから環境変数をエクスポート
eval $(envault export -o -f /path/to/custom.vaulted)
# または
eval $(envault export --output-script-only --file /path/to/custom.vaulted)

# stdinからパスワードを読み込んでエクスポート
echo "password" | envault export -o -p | eval
# または
echo "password" | envault export --output-script-only --password-stdin | eval
```

**注意**: この方法では、`envault export`コマンドはシェルスクリプトを出力するだけで、環境変数を直接設定しません。環境変数を実際に設定するには、上記のように`-o`または`--output-script-only`フラグを使用して、`eval`または`source`コマンドで実行する必要があります。

#### 新しい方法（より簡単）

##### 新しいbashセッションを起動

```bash
# 新しいbashセッションを起動して環境変数を設定
envault export -n
# または
envault export --new-shell

# 特定の暗号化ファイルから環境変数を設定して新しいbashセッションを起動
envault export -n -f /path/to/custom.vaulted
# または
envault export --new-shell --file /path/to/custom.vaulted

# stdinからパスワードを読み込んで新しいbashセッションを起動
echo "password" | envault export -n -p
# または
echo "password" | envault export --new-shell --password-stdin
```

この方法では、envaultが環境変数を設定した新しいbashセッションを起動します。元のシェルには影響を与えませんが、新しいシェル内ですべての環境変数が利用可能になります。

##### 特定のコマンドを実行

```bash
# 環境変数を設定して特定のコマンドを実行
envault export -- node app.js

# 環境変数を設定してdocker-composeを実行
envault export -- docker-compose up

# 環境変数を設定してenvコマンドで確認
envault export -- env | grep SECRET

# 特定の暗号化ファイルから環境変数を設定してコマンドを実行
envault export -f /path/to/custom.vaulted -- python script.py
# または
envault export --file /path/to/custom.vaulted -- python script.py
```

この方法では、envaultが環境変数を設定してから指定されたコマンドを実行します。コマンドとその引数は `--` の後に指定します。

### 環境変数の選択的エクスポート

TUIを使用して、適用する環境変数を選択できます：

```bash
# TUIで環境変数を選択してからエクスポート
envault export -s
# または
envault export --select
# または
envault export select

# TUIで環境変数を選択してから新しいbashセッションを起動
envault export -s -n
# または
envault export select --new-shell

# TUIで環境変数を選択してから特定のコマンドを実行
envault export -s -- npm start
# または
envault export select -- npm start
```

### 環境変数のアンセット

```bash
# .env.vaultedファイルに記載された環境変数をアンセット（重要：evalまたはsourceで実行する必要があります）
eval $(envault unset -o)
# または
eval $(envault unset --output-script-only)
# または
source <(envault unset --output-script-only)

# 特定の暗号化ファイルに記載された環境変数をアンセット
eval $(envault unset -o -f /path/to/custom.vaulted)
# または
eval $(envault unset --output-script-only --file /path/to/custom.vaulted)

# stdinからパスワードを読み込んでアンセット
echo "password" | envault unset -o -p | eval
# または
echo "password" | envault unset --output-script-only --password-stdin | eval

# TUIで環境変数を選択してからアンセット
eval $(envault unset -s -o)
# または
eval $(envault unset --select --output-script-only)
```

**注意**: `envault unset`コマンドもシェルスクリプトを出力するだけで、環境変数を直接アンセットしません。環境変数を実際にアンセットするには、上記のように`-o`または`--output-script-only`フラグを使用して、`eval`または`source`コマンドで実行する必要があります。

### 暗号化されたファイルの内容を確認

```bash
# .env.vaultedファイルの内容を復号化して表示
envault dump

# 特定の暗号化ファイルの内容を表示
envault dump -f /path/to/custom.vaulted
# または
envault dump --file /path/to/custom.vaulted

# 復号化した内容をファイルに保存
envault dump > decrypted.env

# stdinからパスワードを読み込んで復号化
echo "password" | envault dump -p
# または
echo "password" | envault dump --password-stdin
```

### ヘルプとバージョン情報

```bash
# ヘルプを表示
envault help
# 特定のコマンドのヘルプを表示
envault export --help
envault unset --help
envault dump --help

# バージョン情報を表示
envault version
# または
envault --version
```

### コマンド構造

```
envault [オプション] <.envファイル> # 暗号化
envault export [オプション]          # 環境変数のエクスポート
envault export select [オプション]   # 選択的なエクスポート
envault unset [オプション]           # 環境変数のアンセット
envault dump [オプション]            # 暗号化ファイルの内容表示
envault version                     # バージョン表示
envault help                        # ヘルプ表示
```

## コマンドラインパーサー

envaultは、コマンドライン引数の解析に`github.com/spf13/cobra`ライブラリを使用しています。これにより、以下の利点があります：

- ネストされたサブコマンドの自然なサポート
- 自動的な詳細なヘルプメッセージ
- フラグのショートハンドと長形式の自動サポート
- シェル補完の生成（将来的な機能）

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