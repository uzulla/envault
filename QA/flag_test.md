# envault フラグテスト手順書

## 概要

この文書は、envaultコマンドのフラグ解析に関する問題を診断するための手順を提供します。特に`--select`フラグの処理に関する問題に焦点を当てています。

## テスト環境

- macOS または Linux
- Go 1.23.0 以上

## 診断手順

### 1. フラグの認識テスト

#### 1.1 ヘルプ表示での確認

```bash
# ビルドしたenvaultコマンドでフラグが認識されるか確認
./envault -h | grep select
./envault export -h | grep select
```

期待される結果:
- `--select` フラグがヘルプに表示されること

#### 1.2 デバッグ出力での確認

```bash
# selectフラグのみを指定
./envault export --select 2>&1 | head -2

# 短い形式で指定
./envault export -select 2>&1 | head -2

# 値付きで指定
./envault export --select=true 2>&1 | head -2

# 複数フラグの組み合わせ
./envault export --select --new-shell 2>&1 | head -2
```

期待される結果:
- デバッグ出力の1行目に `select=true` が含まれていること

#### 1.3 実行時の動作確認

*注意: このテストでは一時的に.env.vaultedファイルが必要です。テスト用に小さなファイルを暗号化してください*

```bash
# テスト用の.envファイルを作成
echo "TEST=value" > test.env

# 暗号化（パスワードは「test」などシンプルなものを使用）
./envault test.env

# selectフラグでエクスポート
./envault export --select
```

期待される結果:
- パスワード入力後、TUIインターフェースが表示されること

## 問題が発生した場合

フラグが認識されない場合は、以下の情報を収集してください：

1. 使用中のOS情報
```bash
uname -a
```

2. Go言語のバージョン
```bash
go version
```

3. ビルドコマンド
```bash
go build -v -x cmd/envault/main.go
```

4. コマンドラインから実行したコマンドと正確な引数
```
# 例
./envault export --select
```

5. 発生したエラーメッセージの完全な内容

## 追加の診断ツール

問題がフラグの認識や解析にある場合は、以下の診断ツールを使用できます：

```bash
go run debug_cli.go export --select
```

この診断ツールは、フラグの解析過程と結果を詳細に表示します。

## 問題修正のヒント

もし環境特有の問題でフラグが認識されない場合は、以下の方法で問題を回避できる可能性があります：

1. 短いフラグ形式を試す: `--select` の代わりに `-select` を使用
2. 値付きフラグを試す: `--select=true` の形式を使用
3. フラグの順序を変更する: 例えば `--select --new-shell` の代わりに `--new-shell --select` を試す