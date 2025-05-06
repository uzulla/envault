# envault 品質保証 (QA) ガイド

このディレクトリには、envaultのテストと品質保証に関連するツールとドキュメントが含まれています。

## ディレクトリ内容

- `manual_testing.md` - マニュアルテスト手順書
- `UNIT_TEST.md` - ユニットテストの詳細ガイド
- `test_files/` - テスト用ファイルを格納するディレクトリ
  - `test.env` - テスト用の標準環境変数ファイル
  - `multiple_vars.env` - 多数の環境変数を含むテストファイル
  - `test_env_generator.sh` - テスト用の.envファイル生成スクリプト
- `results/` - テスト結果ファイルを格納するディレクトリ

## テスト方法

### 1. 自動テスト (ユニットテスト)

envaultには包括的な自動テストが実装されています。実行方法と詳細については [UNIT_TEST.md](./UNIT_TEST.md) を参照してください。

### 2. マニュアルテスト

自動テストでは検証が難しい部分や、実際のユーザー体験に関わる部分については、マニュアルテストを実施します。

詳細なテスト手順は `manual_testing.md` を参照してください。テスト結果は `results` dirの中に `manual_test_results-YYYYMMDD-HHMMSS.md` 形式のファイル名で保存します。

#### テスト用環境の準備

1. envaultをビルドします：
   ```bash
   cd /path/to/envault
   go build -o envault cmd/envault/main.go
   ```

2. テスト用の.envファイルを生成（必要な場合）：
   ```bash
   # デフォルト設定で.envファイルを生成
   ./QA/test_files/test_env_generator.sh
   
   # カスタム名と環境変数の数を指定
   ./QA/test_files/test_env_generator.sh ./QA/test_files/custom_test.env 20
   ```

#### マニュアルテスト対象機能

1. **基本機能**
   - 暗号化: `envault .env`
   - エクスポート: `envault export`
   - アンセット: `envault unset`
   - ダンプ: `envault dump`
   - ヘルプとバージョン: `envault help`, `envault --version`

2. **TUI選択機能**
   - 環境変数の選択的エクスポート: `envault export -s`
   - 環境変数の選択的アンセット: `envault unset -s`
   - サブコマンド形式: `envault export select`

3. **コマンド実行機能**
   - 環境変数を設定して特定のコマンドを実行: `envault export -- command args`
   - 選択した環境変数でコマンドを実行: `envault export -s -- command args`

4. **シェルセッション機能**
   - 新しいシェルでの環境変数設定: `envault export -n`
   - 選択した変数での新しいシェル: `envault export -s -n`

5. **エラーケース**
   - 存在しないファイルの処理
   - 不正な形式のファイルの処理
   - 誤ったパスワードでの復号化

### 3. セキュリティテスト

envaultのセキュリティ機能を検証するための特別なテストです。

1. **暗号化強度の検証**
   - 暗号化ファイル内に平文が含まれていないことの確認
   - バイナリエディタを使用した暗号化データの検査
   - 暗号化結果のランダム性検証

2. **ファイルセキュリティ**
   - 暗号化ファイルのパーミッション検証（0600）
   - メモリ内パスワード処理の安全性

3. **パスワード強度**
   - Argon2idによる安全な鍵導出の検証
   - 攻撃に対する耐性の評価

## テスト結果の報告

テスト結果は以下の形式で記録します：

1. テスト実施日時
2. テスト環境（OS、Go バージョン）
3. テスト対象のenvaultバージョン
4. 各テスト項目の結果（成功/失敗/スキップ）
5. 発見された問題点や改善提案
6. スクリーンショットや補足情報（必要に応じて）

テスト結果は `results` dirの中に `all_test_results-YYYYMMDD-HHMMSS.md` 形式のファイル名で保存します。

## CI/CD統合（将来の計画）

将来的には、以下のCI/CD統合を計画しています：

1. GitHub Actionsによる自動テスト実行
2. カバレッジレポート生成
3. セキュリティ脆弱性スキャン
4. クロスプラットフォームビルドとテスト
