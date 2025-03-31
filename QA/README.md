# envault QAツール

このディレクトリには、envaultのテストと品質保証に関連するツールとドキュメントが含まれています。

## 内容

- `manual_testing.md` - マニュアルテスト手順書
- `test_env_generator.sh` - テスト用の.envファイル生成スクリプト

## 使用方法

### テスト用.envファイルの生成

テスト用の.envファイルを生成するには、以下のコマンドを実行します：

```bash
# デフォルト設定で.envファイルを生成
./test_env_generator.sh

# カスタム名と環境変数の数を指定
./test_env_generator.sh test.env 10
```

### マニュアルテストの実行

`manual_testing.md`に記載されているテスト手順に従って、envaultの各機能をテストします。テスト結果は同ファイル内のテスト結果記録表に記入してください。

## テスト環境の準備

1. envaultをビルドします：
   ```bash
   cd /path/to/envault
   go build -o envault cmd/envault/main.go
   ```

2. テスト用の.envファイルを生成します：
   ```bash
   ./QA/test_env_generator.sh
   ```

3. マニュアルテスト手順に従ってテストを実行します。
