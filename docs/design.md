# envault CLI ツール設計書

## 概要

envaultは、環境変数を含む.envファイルを安全に暗号化・復号化するためのCLIツールです。このツールは、.envファイルを暗号化して.env.vaultedファイルを作成し、必要に応じて環境変数をエクスポートまたはアンセットする機能を提供します。

## コマンド構造

envaultは以下の主要コマンドをサポートします：

1. `envault encrypt .env` - .envファイルを暗号化して.env.vaultedファイルを作成
2. `envault export` - .env.vaultedファイルから環境変数をエクスポート
   - `envault export select` - TUIを使って環境変数を選択してからエクスポート
3. `envault unset` - .env.vaultedファイルに記載された環境変数をアンセット
4. `envault dump` - .env.vaultedファイルの内容を復号化して表示
5. `envault version` - バージョン情報を表示
6. `envault help` - ヘルプを表示

共通のグローバルオプション：
- `-p, --password-stdin` - パスワードを標準入力から読み込む（対話モードの代わりに）
- `-f, --file` - 使用する.env.vaultedファイルのパスを指定（デフォルトは現在のディレクトリの.env.vaulted）

exportコマンドの追加オプション：
- `-o, --output-script-only` - スクリプトのみを出力（情報メッセージなし）
- `-n, --new-shell` - 新しいbashセッションを起動して環境変数を設定
- `-s, --select` - 適用する環境変数をTUIで選択する

unsetコマンドの追加オプション：
- `-o, --output-script-only` - スクリプトのみを出力（情報メッセージなし）
- `-s, --select` - 適用する環境変数をTUIで選択する

## アーキテクチャ

```
+------------------+     +------------------+     +------------------+
|                  |     |                  |     |                  |
|  CLI Interface   | --> |  Core Logic      | --> |  File Operations |
|  (Cobra)         |     |                  |     |                  |
+------------------+     +------------------+     +------------------+
        |                        |                        |
        v                        v                        v
+------------------+     +------------------+     +------------------+
|                  |     |                  |     |                  |
|  Password Input  |     |  Encryption      |     |  Environment     |
|  Handling        |     |  Operations      |     |  Variable Ops    |
|                  |     |                  |     |                  |
+------------------+     +------------------+     +------------------+
                                                         |
                                                         v
                                               +------------------+
                                               |                  |
                                               |  TUI Interface   |
                                               |  (Bubbletea)     |
                                               |                  |
                                               +------------------+
```

### コンポーネント

1. **CLI Interface**: Cobraライブラリを使用したコマンドライン引数の解析と適切なアクションの呼び出し
2. **Core Logic**: 主要な機能ロジックの実装
3. **File Operations**: ファイルの読み書き操作
4. **Password Input Handling**: 対話式およびstdinからのパスワード入力処理
5. **Encryption Operations**: 暗号化・復号化の処理
6. **Environment Variable Operations**: 環境変数のエクスポート・アンセット処理
7. **TUI Interface**: Bubbleteaライブラリを使用した対話的な環境変数選択インターフェース

## ファイル構造

```
envault/
├── cmd/
│   └── envault/
│       └── main.go           # エントリーポイント
├── internal/
│   ├── cli/
│   │   └── cli.go            # CLIインターフェース（Cobra実装）
│   ├── crypto/
│   │   └── crypto.go         # 暗号化・復号化機能
│   ├── env/
│   │   └── env.go            # 環境変数操作
│   ├── file/
│   │   └── file.go           # ファイル操作
│   └── tui/
│       ├── interface.go      # TUIインターフェース定義
│       ├── factory.go        # TUIファクトリー
│       └── bubbletea_provider.go  # Bubbletea実装
├── pkg/
│   └── utils/
│       └── utils.go          # ユーティリティ関数
└── test/                     # テストファイル
```

## 暗号化仕様

### アルゴリズム選択

AES-256-GCMを使用して、機密性と完全性の両方を確保します。これは現代的で強力な暗号化アルゴリズムであり、認証付き暗号化を提供します。

### 暗号化プロセス

1. ユーザーからパスワードを取得
2. パスワードからArgon2idを使用して鍵を導出（KDF）
3. ランダムなnonce（初期化ベクトル）を生成
4. .envファイルの内容全体を暗号化（キーと値の両方を含む）
5. 暗号化されたデータ、nonce、およびメタデータを.env.vaultedファイルに保存

### ファイル形式

.env.vaultedファイルは以下の形式で保存されます：

```
ENVAULT1       # マジックバイト（ファイル識別子）
[salt]         # 鍵導出用のソルト（16バイト）
[nonce]        # 暗号化用のnonce（12バイト）
[encrypted]    # 暗号化されたデータ
```

この形式により、ファイル内のどの環境変数キーが含まれているかを特定することはできません。

## パスワード処理

### 対話モード

デフォルトでは、ユーザーに対話的にパスワードの入力を求めます。入力されたパスワードはエコーバックされません（セキュリティのため）。

### stdin モード

`--password-stdin`オプションが指定された場合、パスワードは標準入力から読み込まれます。これにより、スクリプト内での使用が可能になります。

## 環境変数の処理

### エクスポート方式

`envault export`コマンドは3つの方式をサポートします：

1. **スクリプト評価方式**: `eval $(envault export -o)` または `source <(envault export --output-script-only)`
2. **新しいシェル方式**: `envault export -n` または `envault export --new-shell`
3. **コマンド実行方式**: `envault export -- node app.js` または `envault export -- docker-compose up`

### 環境変数の選択的適用

`envault export -s` または `envault export select` コマンドを使用すると、TUIインターフェースで環境変数を選択的に適用できます。

### アンセット

`envault unset`コマンドは以下の処理を行います：

1. .env.vaultedファイルを読み込む
2. パスワードを取得して復号化
3. 復号化された環境変数を現在のシェルからアンセット

同様に、`-s` オプションでTUIを使用して選択的にアンセットすることも可能です。

## プラットフォーム互換性

LinuxとmacOSの両方で動作するように設計されています。シングルバイナリとして配布され、追加の依存関係は必要ありません。

## セキュリティ考慮事項

- パスワードはメモリ内で安全に処理され、使用後にゼロクリアされます
- 暗号化されたファイルからは環境変数のキー名を特定できません
- 強力な暗号化アルゴリズム（AES-256-GCM）を使用
- 強力な鍵導出関数（Argon2id）を使用

## コマンドラインインターフェース

### コマンドラインパーサー

コマンドライン引数の解析には、Goの標準ライブラリではなく、`github.com/spf13/cobra`を使用しています。Cobraは以下の利点を提供します：

1. ネストされたサブコマンドのサポート
2. 自動的なヘルプメッセージの生成
3. 引数の検証
4. フラグのエイリアスと自動短縮形
5. シェル補完の生成

### コマンド構造

Cobraを使用してコマンド階層が以下のように構造化されています：

```
envault
├── encrypt
├── export
│   └── select
├── unset
├── dump
└── version
```

各コマンドには適切なフラグと引数が定義されています。