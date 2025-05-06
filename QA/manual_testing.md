# envault マニュアルテスト手順

## 前提条件

- envaultがビルドされていること
- テスト用の.envファイルがQAディレクトリに用意されていること

## テスト環境

- Linux (Ubuntu)
- macOS
- ターミナルはカラー表示対応であること（TUI機能のテスト用）

## テスト項目

### 1. ヘルプとバージョン表示のテスト

#### 1.1 ヘルプ表示

1. 以下のコマンドを実行

   ```bash
   ./envault --help
   # または
   ./envault -h
   # または
   ./envault help
   ```

2. 確認事項:
   - ヘルプメッセージが正しく表示されること
   - コマンド一覧とその説明が表示されること
   - グローバルオプションが表示されること

#### 1.2 サブコマンドのヘルプ表示

1. 以下のコマンドを実行

   ```bash
   ./envault export --help
   ./envault unset --help
   ./envault dump --help
   ```

2. 確認事項:
   - 各サブコマンドの詳細なヘルプが表示されること
   - サブコマンド特有のオプションが表示されること

#### 1.3 バージョン表示

1. 以下のコマンドを実行

   ```bash
   ./envault version
   # または
   ./envault --version
   ```

2. 確認事項:
   - バージョン情報が正しく表示されること
   - 表示形式が `envault version X.X.X` であること

### 2. 暗号化機能のテスト

#### 2.1 基本的な暗号化

1. QAディレクトリのテスト用の.envファイルを使用

   ```env
   TEST_VAR1=value1
   TEST_VAR2=value2
   ```

2. 以下のコマンドを実行

   ```bash
   ./envault QA/test_files/test.env
   ```

3. 確認事項:
   - QA/.env.vaultedファイルが作成されること
   - パスワードの入力が求められること
   - パスワード確認（再入力）が求められること
   - 処理が正常に完了すること

#### 2.2 既存の.env.vaultedファイルの上書き確認

1. 既にQA/.env.vaultedファイルが存在する状態で暗号化を実行
2. 確認事項:
   - 上書きの確認メッセージが表示されること
   - 「Y」を入力すると上書きされること
   - 「N」を入力すると処理がキャンセルされること

#### 2.3 カスタム出力パスの指定

1. 以下のコマンドを実行

   ```bash
   ./envault QA/test.env -f QA/custom.vaulted
   # または
   ./envault QA/test.env --file QA/custom.vaulted
   ```

2. 確認事項:
   - 指定したパスにファイルが作成されること
   - 処理が正常に完了すること

#### 2.4 stdinからのパスワード入力

1. 以下のコマンドを実行

   ```bash
   echo "password" | ./envault QA/test.env -p
   # または
   echo "password" | ./envault QA/test.env --password-stdin
   ```

2. 確認事項:
   - パスワード入力がスキップされること
   - .env.vaultedファイルが正常に作成されること

### 3. エクスポート機能のテスト

#### 3.1 基本的なエクスポート（従来の方法）

1. 暗号化されたQA/.env.vaultedファイルがある状態で以下のコマンドを実行

   ```bash
   eval $(./envault export -f QA/.env.vaulted -o)
   # または
   eval $(./envault export --file QA/.env.vaulted --output-script-only)
   # または
   source <(./envault export -f QA/.env.vaulted -o)
   ```

2. 確認事項:
   - パスワードの入力が求められること
   - 正しいパスワードを入力すると環境変数がエクスポートされること
   - `echo $TEST_VAR1`で値が表示されること

#### 3.2 誤ったパスワードでのエクスポート

1. 誤ったパスワードを入力
2. 確認事項:
   - エラーメッセージが表示されること
   - 環境変数がエクスポートされないこと

#### 3.3 stdinからのパスワード入力

1. 以下のコマンドを実行

   ```bash
   echo "password" | ./envault export -f QA/.env.vaulted -p -o | eval
   # または
   echo "password" | ./envault export --file QA/.env.vaulted --password-stdin --output-script-only | eval
   ```

2. 確認事項:
   - 環境変数が正常にエクスポートされること

#### 3.4 TUIを使用した選択的エクスポート (フラグ形式)

1. 暗号化されたQA/.env.vaultedファイルがある状態で以下のコマンドを実行

   ```bash
   eval $(./envault export -f QA/.env.vaulted -s -o)
   # または
   eval $(./envault export --file QA/.env.vaulted --select --output-script-only)
   ```

2. 確認事項:
   - パスワードの入力が求められること
   - TUIインターフェースが表示されること
   - 上下キーで選択項目を移動できること
   - スペースキーで環境変数のON/OFFを切り替えられること
   - Enterキーで選択を確定できること
   - 選択した環境変数のみがエクスポートされ、選択しなかった環境変数はエクスポートされないこと
   - コメント付きの環境変数は、コメントがTUI上で表示されること

#### 3.5 TUIを使用した選択的エクスポート (サブコマンド形式)

1. 暗号化されたQA/.env.vaultedファイルがある状態で以下のコマンドを実行

   ```bash
   eval $(./envault export select -f QA/.env.vaulted -o)
   # または
   eval $(./envault export select --file QA/.env.vaulted --output-script-only)
   ```

2. 確認事項:
   - パスワードの入力が求められること
   - TUIインターフェースが表示されること
   - 上下キーで選択項目を移動できること
   - スペースキーで環境変数のON/OFFを切り替えられること
   - Enterキーで選択を確定できること
   - 選択した環境変数のみがエクスポートされること

#### 3.6 TUIを使用した選択的エクスポート（キャンセル）

1. TUIインターフェースが表示された状態で以下のキーを押す
   - qキー
   - Escキー

2. 確認事項:
   - TUIがキャンセルされること
   - キャンセル後の動作が適切であること

#### 3.7 新しいbashセッションの起動

1. 暗号化されたQA/.env.vaultedファイルがある状態で以下のコマンドを実行

   ```bash
   ./envault export -f QA/.env.vaulted -n
   # または
   ./envault export --file QA/.env.vaulted --new-shell
   ```

2. 確認事項:
   - パスワードの入力が求められること
   - 正しいパスワードを入力すると新しいbashセッションが起動すること
   - 新しいシェル内で`echo $TEST_VAR1`を実行すると値が表示されること
   - 新しいシェルを終了（`exit`コマンド）すると、元のシェルでは環境変数が設定されていないこと

#### 3.8 TUIで選択して新しいbashセッションを起動

1. 暗号化されたQA/.env.vaultedファイルがある状態で以下のコマンドを実行

   ```bash
   ./envault export -f QA/.env.vaulted -s -n
   # または
   ./envault export --file QA/.env.vaulted --select --new-shell
   # または
   ./envault export select -f QA/.env.vaulted -n
   ```

2. 確認事項:
   - パスワードの入力が求められること
   - TUIインターフェースが表示されること
   - 一部の環境変数を選択解除してEnterキーで確定すると、選択した環境変数のみが設定された新しいbashセッションが起動すること
   - 選択解除した環境変数は新しいシェルでは設定されていないこと

#### 3.9 コマンド実行オプション

1. 暗号化されたQA/.env.vaultedファイルがある状態で以下のコマンドを実行

   ```bash
   ./envault export -f QA/.env.vaulted -- env | grep TEST_VAR
   # または
   ./envault export --file QA/.env.vaulted -- env | grep TEST_VAR
   ```

2. 確認事項:
   - パスワードの入力が求められること
   - 正しいパスワードを入力するとenvコマンドが実行され、TEST_VARで始まる環境変数が表示されること
   - コマンド実行後、元のシェルでは環境変数が設定されていないこと

#### 3.10 TUIで選択してコマンド実行

1. 以下のコマンドを実行

   ```bash
   ./envault export -f QA/.env.vaulted -s -- env | grep TEST_VAR
   # または
   ./envault export --file QA/.env.vaulted --select -- env | grep TEST_VAR
   # または
   ./envault export select -f QA/.env.vaulted -- env | grep TEST_VAR
   ```

2. 確認事項:
   - TUIインターフェースが表示されること
   - 選択した環境変数のみが設定された状態でコマンドが実行されること
   - 選択解除した環境変数はコマンド実行時の環境には含まれないこと

#### 3.11 複雑なコマンド実行

1. 以下のコマンドを実行

   ```bash
   ./envault export -f QA/.env.vaulted -- bash -c "echo TEST_VAR1=$TEST_VAR1 TEST_VAR2=$TEST_VAR2"
   # または
   ./envault export --file QA/.env.vaulted -- bash -c "echo TEST_VAR1=$TEST_VAR1 TEST_VAR2=$TEST_VAR2"
   ```

2. 確認事項:
   - 指定したbashコマンドが実行され、環境変数の値が正しく表示されること

### 4. アンセット機能のテスト

#### 4.1 基本的なアンセット

1. 環境変数がエクスポートされた状態で以下のコマンドを実行

   ```bash
   eval $(./envault unset -f QA/.env.vaulted -o)
   # または
   eval $(./envault unset --file QA/.env.vaulted --output-script-only)
   # または
   source <(./envault unset -f QA/.env.vaulted -o)
   ```

2. 確認事項:
   - パスワードの入力が求められること
   - 正しいパスワードを入力すると環境変数がアンセットされること
   - `echo $TEST_VAR1`で値が表示されないこと

#### 4.2 TUIを使用した選択的アンセット

1. 環境変数がエクスポートされた状態で以下のコマンドを実行

   ```bash
   eval $(./envault unset -f QA/.env.vaulted -s -o)
   # または
   eval $(./envault unset --file QA/.env.vaulted --select --output-script-only)
   ```

2. 確認事項:
   - パスワードの入力が求められること
   - TUIインターフェースが表示されること
   - 上下キーで選択項目を移動できること
   - スペースキーで環境変数のON/OFFを切り替えられること
   - Enterキーで選択を確定できること
   - 選択した環境変数のみがアンセットされ、選択しなかった環境変数はアンセットされないこと

#### 4.3 stdinからのパスワード入力

1. 以下のコマンドを実行

   ```bash
   echo "password" | eval $(./envault unset -f QA/.env.vaulted -p -o)
   # または
   echo "password" | eval $(./envault unset --file QA/.env.vaulted --password-stdin --output-script-only)
   ```

2. 確認事項:
   - 環境変数が正常にアンセットされること

### 5. ダンプ機能のテスト

#### 5.1 基本的なダンプ

1. 暗号化されたQA/.env.vaultedファイルがある状態で以下のコマンドを実行

   ```bash
   ./envault dump -f QA/.env.vaulted
   # または
   ./envault dump --file QA/.env.vaulted
   ```

2. 確認事項:
   - パスワードの入力が求められること
   - 正しいパスワードを入力すると.envファイルの内容が表示されること
   - 表示された内容が元のQA/test.envファイルと一致すること

#### 5.2 カスタムファイルのダンプ

1. 以下のコマンドを実行

   ```bash
   ./envault dump -f QA/custom.vaulted
   # または
   ./envault dump --file QA/custom.vaulted
   ```

2. 確認事項:
   - 指定したファイルの内容が正常に表示されること

#### 5.3 出力のリダイレクト

1. 以下のコマンドを実行

   ```bash
   ./envault dump -f QA/.env.vaulted > QA/decrypted.env
   # または
   ./envault dump --file QA/.env.vaulted > QA/decrypted.env
   ```

2. 確認事項:
   - QA/decrypted.envファイルが作成されること
   - ファイルの内容が元のQA/test.envファイルと一致すること

#### 5.4 stdinからのパスワード入力

1. 以下のコマンドを実行

   ```bash
   echo "password" | ./envault dump -f QA/.env.vaulted -p
   # または
   echo "password" | ./envault dump --file QA/.env.vaulted --password-stdin
   ```

2. 確認事項:
   - .envファイルの内容が正常に表示されること

### 6. TUI機能のテスト

#### 6.1 TUIインターフェースの応答性

1. 以下のコマンドを実行

   ```bash
   ./envault export -f QA/.env.vaulted -s -n
   # または
   ./envault export --file QA/.env.vaulted --select --new-shell
   # または
   ./envault export select -f QA/.env.vaulted -n
   ```

2. 確認事項:
   - TUIの表示が適切に描画されること
   - キーボード操作に対してレスポンシブであること
   - ヘルプ表示が画面下部に表示されること
   - ターミナルサイズを変更した場合、TUIが適切にリサイズされること

#### 6.2 複数の環境変数を含むファイルでのTUI表示

1. 多数の環境変数を含むテストファイルを暗号化後、エクスポートを実行

   ```bash
   # まず暗号化
   ./envault QA/test_files/multiple_vars.env -f QA/test_files/multiple_vars.env.vaulted
   
   # そしてエクスポート
   ./envault export -f QA/test_files/multiple_vars.env.vaulted -s
   # または
   ./envault export --file QA/test_files/multiple_vars.env.vaulted --select
   # または
   ./envault export select -f QA/test_files/multiple_vars.env.vaulted
   ```

2. 確認事項:
   - 多数の環境変数が適切に表示されること
   - スクロールが必要な場合、適切にスクロール表示されること
   - 長いキー名や値を持つ環境変数が適切に表示されること

#### 6.3 TUIのキャンセル操作

1. TUIモードでエクスポートを実行し、`q`キーまたは`Esc`キーでキャンセル

   ```bash
   ./envault export -f QA/.env.vaulted -s
   # または
   ./envault export --file QA/.env.vaulted --select
   # または
   ./envault export select -f QA/.env.vaulted
   ```

2. 確認事項:
   - `q`キーまたは`Esc`キーでTUIがキャンセルされること
   - キャンセル時に適切な処理が行われること

#### 6.4 TUIでのコメント表示

1. コメント付きの環境変数ファイルを作成し、TUIモードでエクスポートを実行

   ```bash
   cat > QA/test_files/comment_test.env << EOF
   # This is a test variable
   TEST_COMMENT_VAR=value

   # Another test variable with description
   # This spans multiple lines
   TEST_MULTILINE_COMMENT=value
   EOF

   ./envault QA/test_files/comment_test.env -f QA/test_files/comment_test.env.vaulted
   ./envault export -f QA/test_files/comment_test.env.vaulted -s
   ```

2. 確認事項:
   - コメントがTUI上で適切に表示されること
   - 複数行のコメントが連結されて表示されること

### 7. エラーケースのテスト

#### 7.1 存在しない.envファイルの暗号化

1. 存在しない.envファイルを指定して暗号化を実行

   ```bash
   ./envault QA/nonexistent.env
   ```

2. 確認事項:
   - 適切なエラーメッセージが表示されること

#### 7.2 存在しない.env.vaultedファイルのエクスポート

1. .env.vaultedファイルが存在しない状態でエクスポートを実行

   ```bash
   ./envault export -f QA/nonexistent.env.vaulted
   # または
   ./envault export --file QA/nonexistent.env.vaulted
   ```

2. 確認事項:
   - 適切なエラーメッセージが表示されること

#### 7.3 不正な形式の.env.vaultedファイルのエクスポート

1. 不正な形式の.env.vaultedファイルを用意

   ```bash
   echo "invalid data" > QA/test_files/invalid.env.vaulted
   ```

2. エクスポートを実行

   ```bash
   ./envault export -f QA/test_files/invalid.env.vaulted
   # または
   ./envault export --file QA/test_files/invalid.env.vaulted
   ```

3. 確認事項:
   - 適切なエラーメッセージが表示されること

#### 7.4 不正なコマンドラインオプション

1. 不正なオプションでコマンドを実行

   ```bash
   ./envault --invalid-option
   ./envault export --invalid-option
   ```

2. 確認事項:
   - 適切なエラーメッセージが表示されること
   - ヘルプメッセージへの案内があること

### 8. Cobraコマンド構造のテスト

#### 8.1 ヘルプメッセージの整合性

1. 各コマンドのヘルプを表示

   ```bash
   ./envault --help
   ./envault export --help
   ./envault export select --help
   ./envault unset --help
   ./envault dump --help
   ./envault version --help
   ```

2. 確認事項:
   - すべてのコマンド、サブコマンドのヘルプメッセージが適切に表示されること
   - コマンドの説明が正確であること
   - オプションの説明が正確であること
   - 使用例が含まれていること

#### 8.2 コマンド補完のテスト

1. 部分的なコマンド名を入力してタブ補完を試みる
   
   ```bash
   ./envault ex[TAB]
   ./envault uns[TAB]
   ./envault export sel[TAB]
   ```

2. 確認事項:
   - タブ補完が適切に動作すること（ターミナルがタブ補完をサポートしている場合）

## テスト結果記録

テスト結果は `QA/test_results-YYYYMMDD-HHMMSS.md` ファイルに記録してください。

| テスト項目 | 期待結果 | 実際の結果 | 合否 | 備考 |
|------------|----------|------------|------|------|
| 1.1        |          |            |      |      |
| 1.2        |          |            |      |      |
| 1.3        |          |            |      |      |
| 2.1        |          |            |      |      |
| ...        |          |            |      |      |
| 8.2        |          |            |      |      |

## テスト後のクリーンアップ

テスト完了後、以下のコマンドを実行して作成した一時ファイルをクリーンアップしてください。

```bash
# テスト中に作成した暗号化ファイルや一時ファイルを削除
rm -f QA/test_files/.env.vaulted
rm -f QA/test_files/custom.vaulted
rm -f QA/test_files/multiple_vars.env.vaulted
rm -f QA/test_files/comment_test.env
rm -f QA/test_files/comment_test.env.vaulted
rm -f QA/test_files/invalid.env.vaulted
rm -f QA/test_files/decrypted.env
rm -f QA/test_files/generated_test.env
rm -f QA/test_files/test2.env
rm -f QA/test_files/test2.env.vaulted

# QA/test_files ディレクトリ内の設定されていないテストファイルのみが残るようにします
ls -la QA/test_files/
```