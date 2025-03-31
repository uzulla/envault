# envault マニュアルテスト手順

## 前提条件
- envaultがビルドされていること
- テスト用の.envファイルがQAディレクトリに用意されていること

## テスト環境
- Linux (Ubuntu)
- macOS

## テスト項目

### 1. 暗号化機能のテスト

#### 1.1 基本的な暗号化
1. QAディレクトリのテスト用の.envファイルを使用
   ```
   TEST_VAR1=value1
   TEST_VAR2=value2
   ```
2. 以下のコマンドを実行
   ```
   ./envault QA/test.env
   ```
3. 確認事項:
   - QA/.env.vaultedファイルが作成されること
   - パスワードの入力が求められること
   - 処理が正常に完了すること

#### 1.2 既存の.env.vaultedファイルの上書き確認
1. 既にQA/.env.vaultedファイルが存在する状態で暗号化を実行
2. 確認事項:
   - 上書きの確認メッセージが表示されること
   - 「Y」を入力すると上書きされること
   - 「N」を入力すると処理がキャンセルされること

### 2. エクスポート機能のテスト

#### 2.1 基本的なエクスポート（従来の方法）
1. 暗号化されたQA/.env.vaultedファイルがある状態で以下のコマンドを実行
   ```
   eval $(./envault export --file QA/.env.vaulted --output-script-only)
   ```
2. 確認事項:
   - パスワードの入力が求められること
   - 正しいパスワードを入力すると環境変数がエクスポートされること
   - `echo $TEST_VAR1`で値が表示されること

#### 2.2 誤ったパスワードでのエクスポート
1. 誤ったパスワードを入力
2. 確認事項:
   - エラーメッセージが表示されること
   - 環境変数がエクスポートされないこと

#### 2.3 stdinからのパスワード入力
1. 以下のコマンドを実行
   ```
   echo "password" | eval $(./envault export --file QA/.env.vaulted --password-stdin --output-script-only)
   ```
2. 確認事項:
   - 環境変数が正常にエクスポートされること

#### 2.4 新しいbashセッションの起動
1. 暗号化されたQA/.env.vaultedファイルがある状態で以下のコマンドを実行
   ```
   ./envault export --file QA/.env.vaulted --new-shell
   ```
2. 確認事項:
   - パスワードの入力が求められること
   - 正しいパスワードを入力すると新しいbashセッションが起動すること
   - 新しいシェル内で`echo $TEST_VAR1`を実行すると値が表示されること
   - 新しいシェルを終了（`exit`コマンド）すると、元のシェルでは環境変数が設定されていないこと

#### 2.5 コマンド実行オプション
1. 暗号化されたQA/.env.vaultedファイルがある状態で以下のコマンドを実行
   ```
   ./envault export --file QA/.env.vaulted -- env | grep TEST_VAR
   ```
2. 確認事項:
   - パスワードの入力が求められること
   - 正しいパスワードを入力するとenvコマンドが実行され、TEST_VARで始まる環境変数が表示されること
   - コマンド実行後、元のシェルでは環境変数が設定されていないこと

#### 2.6 複雑なコマンド実行
1. 以下のコマンドを実行
   ```
   ./envault export --file QA/.env.vaulted -- bash -c "echo TEST_VAR1=$TEST_VAR1 TEST_VAR2=$TEST_VAR2"
   ```
2. 確認事項:
   - 指定したbashコマンドが実行され、環境変数の値が正しく表示されること

### 3. アンセット機能のテスト

#### 3.1 基本的なアンセット
1. 環境変数がエクスポートされた状態で以下のコマンドを実行
   ```
   ./envault unset --file QA/.env.vaulted
   ```
2. 確認事項:
   - パスワードの入力が求められること
   - 正しいパスワードを入力すると環境変数がアンセットされること
   - `echo $TEST_VAR1`で値が表示されないこと

#### 3.2 stdinからのパスワード入力
1. 以下のコマンドを実行
   ```
   echo "password" | ./envault unset --file QA/.env.vaulted --password-stdin
   ```
2. 確認事項:
   - 環境変数が正常にアンセットされること

### 4. ダンプ機能のテスト

#### 4.1 基本的なダンプ
1. 暗号化されたQA/.env.vaultedファイルがある状態で以下のコマンドを実行
   ```
   ./envault dump --file QA/.env.vaulted
   ```
2. 確認事項:
   - パスワードの入力が求められること
   - 正しいパスワードを入力すると.envファイルの内容が表示されること
   - 表示された内容が元のQA/test.envファイルと一致すること

#### 4.2 カスタムファイルのダンプ
1. 以下のコマンドを実行
   ```
   ./envault dump --file QA/.env.vaulted
   ```
2. 確認事項:
   - 指定したファイルの内容が正常に表示されること

#### 4.3 出力のリダイレクト
1. 以下のコマンドを実行
   ```
   ./envault dump --file QA/.env.vaulted > QA/decrypted.env
   ```
2. 確認事項:
   - QA/decrypted.envファイルが作成されること
   - ファイルの内容が元のQA/test.envファイルと一致すること

#### 4.4 stdinからのパスワード入力
1. 以下のコマンドを実行
   ```
   echo "password" | ./envault dump --file QA/.env.vaulted --password-stdin
   ```
2. 確認事項:
   - .envファイルの内容が正常に表示されること

### 5. エラーケースのテスト

#### 5.1 存在しない.envファイルの暗号化
1. 存在しない.envファイルを指定して暗号化を実行
   ```
   ./envault QA/nonexistent.env
   ```
2. 確認事項:
   - 適切なエラーメッセージが表示されること

#### 5.2 存在しない.env.vaultedファイルのエクスポート
1. .env.vaultedファイルが存在しない状態でエクスポートを実行
   ```
   ./envault export --file QA/nonexistent.env.vaulted
   ```
2. 確認事項:
   - 適切なエラーメッセージが表示されること

#### 5.3 不正な形式の.env.vaultedファイルのエクスポート
1. 不正な形式の.env.vaultedファイルを用意
   ```
   echo "invalid data" > QA/invalid.env.vaulted
   ```
2. エクスポートを実行
   ```
   ./envault export --file QA/invalid.env.vaulted
   ```
3. 確認事項:
   - 適切なエラーメッセージが表示されること

## テスト結果記録

テスト結果は `QA/test_results.md` ファイルに記録してください。

| テスト項目 | 期待結果 | 実際の結果 | 合否 | 備考 |
|------------|----------|------------|------|------|
| 1.1        |          |            |      |      |
| 1.2        |          |            |      |      |
| 2.1        |          |            |      |      |
| 2.2        |          |            |      |      |
| 2.3        |          |            |      |      |
| 3.1        |          |            |      |      |
| 3.2        |          |            |      |      |
| 4.1        |          |            |      |      |
| 4.2        |          |            |      |      |
| 4.3        |          |            |      |      |
| 4.4        |          |            |      |      |
| 5.1        |          |            |      |      |
| 5.2        |          |            |      |      |
| 5.3        |          |            |      |      |
