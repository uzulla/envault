package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/uzulla/envault/internal/crypto"
	"github.com/uzulla/envault/internal/env"
	"github.com/uzulla/envault/internal/file"
	"github.com/uzulla/envault/internal/tui"
	"github.com/uzulla/envault/pkg/utils"
)

const (
	Version = "0.2.0"
)

var (
	ErrInvalidCommand = errors.New("無効なコマンドです")
)

// コマンドモード
type CommandMode int

const (
	ExportMode CommandMode = iota
	UnsetMode
)

type CLI struct {
	rootCmd       *cobra.Command
	passwordStdin bool
	vaultedFile   string
	newShell      bool
	selectVars    bool // 環境変数を選択するオプション
}

func NewCLI() *CLI {
	cli := &CLI{}
	cli.setupCommands()
	return cli
}

func (c *CLI) setupCommands() {
	// ルートコマンド
	c.rootCmd = &cobra.Command{
		Use:     "envault [command]",
		Short:   "環境変数を暗号化して管理するツール",
		Version: Version,
		Long: `envault は環境変数を安全に管理するためのツールです。
.env ファイルを暗号化し、必要に応じて環境変数をエクスポート/アンセットします。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// ヘルプを表示
			return cmd.Help()
		},
		SilenceUsage: true,
	}

	// 共通フラグ
	c.rootCmd.PersistentFlags().BoolVarP(&c.passwordStdin, "password-stdin", "p", false, "stdinからパスワードを読み込む")
	c.rootCmd.PersistentFlags().StringVarP(&c.vaultedFile, "file", "f", "", "使用する.env.vaultedファイルのパス")
	
	// encrypt コマンド
	encryptCmd := &cobra.Command{
		Use:   "encrypt [オプション] <.envファイル>",
		Short: ".envファイルを暗号化して.env.vaultedファイルを作成",
		Long: `.envファイルを暗号化して.env.vaultedファイルを作成します。
- 基本的な暗号化: envault encrypt .env
- カスタム出力パス: envault encrypt .env -f custom.vaulted`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// 第1引数は .env ファイルパス
			return c.runEncrypt(args[0])
		},
	}
	c.rootCmd.AddCommand(encryptCmd)

	// export コマンド
	exportCmd := &cobra.Command{
		Use:   "export [オプション] [-- <コマンド> [引数...]]",
		Short: ".env.vaultedファイルから環境変数をエクスポート",
		Long: `暗号化された .env.vaulted ファイルから環境変数をエクスポートします。
- スクリプト評価方式: eval $(envault export -o)
- 新しいシェル方式: envault export -n
- コマンド実行方式: envault export -- node app.js`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var cmdArgs []string
			if cmd.ArgsLenAtDash() != -1 {
				cmdArgs = args[cmd.ArgsLenAtDash():]
			}

			outputScriptOnly, _ := cmd.Flags().GetBool("output-script-only")
			return c.runWithVaultedFile(ExportMode, outputScriptOnly, cmdArgs)
		},
	}

	exportCmd.Flags().BoolVarP(&c.selectVars, "select", "s", false, "適用する環境変数をTUIで選択する")
	exportCmd.Flags().BoolP("output-script-only", "o", false, "スクリプトのみを出力（情報メッセージなし）")
	exportCmd.Flags().BoolVarP(&c.newShell, "new-shell", "n", false, "新しいbashセッションを起動して環境変数を設定")
	c.rootCmd.AddCommand(exportCmd)

	// select サブコマンド
	selectCmd := &cobra.Command{
		Use:   "select [オプション] [-- <コマンド> [引数...]]",
		Short: "環境変数を選択してからエクスポート",
		Long:  "TUIインターフェースを使用して環境変数を選択的にエクスポートします。",
		RunE: func(cmd *cobra.Command, args []string) error {
			c.selectVars = true
			var cmdArgs []string
			if cmd.ArgsLenAtDash() != -1 {
				cmdArgs = args[cmd.ArgsLenAtDash():]
			}

			outputScriptOnly, _ := cmd.Flags().GetBool("output-script-only")
			return c.runWithVaultedFile(ExportMode, outputScriptOnly, cmdArgs)
		},
	}

	selectCmd.Flags().BoolP("output-script-only", "o", false, "スクリプトのみを出力（情報メッセージなし）")
	selectCmd.Flags().BoolVarP(&c.newShell, "new-shell", "n", false, "新しいbashセッションを起動して環境変数を設定")
	exportCmd.AddCommand(selectCmd)

	// unset コマンド
	unsetCmd := &cobra.Command{
		Use:   "unset [オプション]",
		Short: ".env.vaultedファイルに記載された環境変数をアンセット",
		Long: `.env.vaulted ファイルに記載された環境変数をアンセットします。
- スクリプト評価方式: eval $(envault unset -o)
- source方式: source <(envault unset --output-script-only)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputScriptOnly, _ := cmd.Flags().GetBool("output-script-only")
			return c.runWithVaultedFile(UnsetMode, outputScriptOnly, nil)
		},
	}

	unsetCmd.Flags().BoolVarP(&c.selectVars, "select", "s", false, "適用する環境変数をTUIで選択する")
	unsetCmd.Flags().BoolP("output-script-only", "o", false, "スクリプトのみを出力（情報メッセージなし）")
	c.rootCmd.AddCommand(unsetCmd)

	// dump コマンド
	dumpCmd := &cobra.Command{
		Use:   "dump [オプション]",
		Short: ".env.vaultedファイルを復号化して内容を表示",
		Long: `.env.vaulted ファイルを復号化して内容を表示します。
復号化した内容をリダイレクトしてファイルに保存することもできます: envault dump > decrypted.env`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runDump()
		},
	}
	c.rootCmd.AddCommand(dumpCmd)

	// version コマンド
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "バージョン情報を表示",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("envault version %s\n", Version)
		},
	}
	c.rootCmd.AddCommand(versionCmd)

	// ヘルプとバージョン情報は cobra が自動的に処理
}

func (c *CLI) Run(args []string) error {
	c.rootCmd.SetArgs(args)
	return c.rootCmd.Execute()
}

func (c *CLI) runEncrypt(envFilePath string) error {
	data, err := file.ReadEnvFile(envFilePath)
	if err != nil {
		return fmt.Errorf(".envファイルの読み込みに失敗しました: %w", err)
	}

	var password string
	if c.passwordStdin {
		password, err = utils.GetPasswordFromStdin()
	} else {
		password, err = utils.GetPasswordInteractive("暗号化用パスワードを入力してください: ")
		if err == nil {
			confirmPassword, confirmErr := utils.GetPasswordInteractive("パスワードを再入力してください: ")
			if confirmErr != nil {
				return confirmErr
			}
			if password != confirmPassword {
				return errors.New("パスワードが一致しません")
			}
		}
	}
	if err != nil {
		return err
	}

	encryptedData, err := crypto.Encrypt(data, password)
	if err != nil {
		return fmt.Errorf("暗号化に失敗しました: %w", err)
	}

	outputPath := c.vaultedFile
	if outputPath == "" {
		dir := filepath.Dir(envFilePath)
		if dir == "." {
			outputPath = file.DefaultVaultedFileName
		} else {
			outputPath = filepath.Join(dir, file.DefaultVaultedFileName)
		}
	}

	if err := file.WriteVaultedFile(encryptedData, outputPath); err != nil {
		return fmt.Errorf("暗号化ファイルの書き込みに失敗しました: %w", err)
	}

	fmt.Printf("暗号化されたファイルを作成しました: %s\n", outputPath)
	return nil
}

// 暗号化ファイルを使用する共通ロジック
func (c *CLI) runWithVaultedFile(mode CommandMode, outputScriptOnly bool, cmdArgs []string) error {
	// デバッグ出力は環境変数で制御
	if os.Getenv("ENVAULT_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "[debug] mode=%v select=%v new-shell=%v file=%s\n", 
			mode, c.selectVars, c.newShell, c.vaultedFile)
	}

	// 暗号化ファイルの読み込みと復号化
	data, err := file.ReadVaultedFile(c.vaultedFile)
	if err != nil {
		return fmt.Errorf(".env.vaultedファイルの読み込みに失敗しました: %w", err)
	}

	var password string
	if c.passwordStdin {
		password, err = utils.GetPasswordFromStdin()
	} else {
		password, err = utils.GetPasswordInteractive("復号化用パスワードを入力してください: ")
	}
	if err != nil {
		return err
	}

	decryptedData, err := crypto.Decrypt(data, password)
	if err != nil {
		return fmt.Errorf("復号化に失敗しました: %w", err)
	}

	// 処理モードによって動作を変更
	if c.selectVars {
		return c.processWithTUI(mode, decryptedData, outputScriptOnly, cmdArgs)
	} else {
		return c.processWithoutTUI(mode, decryptedData, outputScriptOnly, cmdArgs)
	}
}

// TUIを使用して環境変数を処理
func (c *CLI) processWithTUI(mode CommandMode, decryptedData []byte, outputScriptOnly bool, cmdArgs []string) error {
	// コメント付きで環境変数を解析
	envVarList, err := env.ParseEnvContentWithComments(decryptedData)
	if err != nil {
		return fmt.Errorf("環境変数の解析に失敗しました: %w", err)
	}

	// TUIで環境変数を選択
	selectedEnvVars, err := c.selectEnvironmentVariables(envVarList)
	if err != nil {
		return fmt.Errorf("環境変数の選択に失敗しました: %w", err)
	}

	// 有効な環境変数のみをマップに変換
	envVars := env.FilterEnabledEnvVars(selectedEnvVars)

	// 有効な環境変数の数をマップの長さから取得
	enabledCount := len(envVars)
	
	// キャンセルされた場合（有効な環境変数が0の場合）は処理を中止
	if enabledCount == 0 {
		fmt.Fprintf(os.Stderr, "操作がキャンセルされました\n")
		return nil
	}

	// モードに応じた処理
	switch mode {
	case ExportMode:
		if c.newShell {
			return c.runNewShell(envVars, enabledCount)
		} else if len(cmdArgs) > 0 {
			return c.runCommand(envVars, cmdArgs, enabledCount)
		} else if outputScriptOnly {
			script := env.GenerateExportScriptFromEnvVarList(selectedEnvVars)
			fmt.Print(script)
		} else {
			script := env.GenerateExportScriptFromEnvVarList(selectedEnvVars)
			if err := utils.ExecuteScript(script); err != nil {
				return fmt.Errorf("環境変数のエクスポートに失敗しました: %w", err)
			}
			fmt.Fprintf(os.Stderr, "%d個の環境変数をエクスポートしました\n", enabledCount)
		}
	case UnsetMode:
		script := env.GenerateUnsetScriptFromEnvVarList(selectedEnvVars)
		if outputScriptOnly {
			fmt.Print(script)
		} else {
			if err := utils.ExecuteScript(script); err != nil {
				return fmt.Errorf("環境変数のアンセットに失敗しました: %w", err)
			}
			fmt.Fprintf(os.Stderr, "%d個の環境変数をアンセットしました\n", enabledCount)
		}
	}

	return nil
}

// TUIを使用せずに環境変数を処理
func (c *CLI) processWithoutTUI(mode CommandMode, decryptedData []byte, outputScriptOnly bool, cmdArgs []string) error {
	// 従来の方法で環境変数を解析
	envVars, err := env.ParseEnvContent(decryptedData)
	if err != nil {
		return fmt.Errorf("環境変数の解析に失敗しました: %w", err)
	}

	envVarCount := len(envVars)

	// モードに応じた処理
	switch mode {
	case ExportMode:
		if c.newShell {
			return c.runNewShell(envVars, envVarCount)
		} else if len(cmdArgs) > 0 {
			return c.runCommand(envVars, cmdArgs, envVarCount)
		} else if outputScriptOnly {
			script := env.GenerateExportScript(envVars)
			fmt.Print(script)
		} else {
			script := env.GenerateExportScript(envVars)
			if err := utils.ExecuteScript(script); err != nil {
				return fmt.Errorf("環境変数のエクスポートに失敗しました: %w", err)
			}
			fmt.Fprintf(os.Stderr, "%d個の環境変数をエクスポートしました\n", envVarCount)
		}
	case UnsetMode:
		script := env.GenerateUnsetScript(envVars)
		if outputScriptOnly {
			fmt.Print(script)
		} else {
			if err := utils.ExecuteScript(script); err != nil {
				return fmt.Errorf("環境変数のアンセットに失敗しました: %w", err)
			}
			fmt.Fprintf(os.Stderr, "%d個の環境変数をアンセットしました\n", envVarCount)
		}
	}

	return nil
}

// 新しいシェルセッションを起動
func (c *CLI) runNewShell(envVars map[string]string, count int) error {
	// 親プロセスの環境変数に影響を与えないために、子プロセス用の環境変数のみを設定
	envSlice := os.Environ()
	for k, v := range envVars {
		prefix := k + "="
		// 既に存在するエントリを除外
		filtered := make([]string, 0, len(envSlice))
		for _, e := range envSlice {
			if !strings.HasPrefix(e, prefix) {
				filtered = append(filtered, e)
			}
		}
		// 新しい値を追加
		envSlice = filtered
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}
	
	fmt.Fprintf(os.Stderr, "%d個の環境変数を設定して新しいbashセッションを起動します\n", count)
	
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = envSlice
	
	return cmd.Run()
}

// 特定のコマンドを実行
func (c *CLI) runCommand(envVars map[string]string, cmdArgs []string, count int) error {
	// 親プロセスの環境変数に影響を与えないために、子プロセス用の環境変数のみを設定
	envSlice := os.Environ()
	for k, v := range envVars {
		prefix := k + "="
		// 既に存在するエントリを除外
		filtered := make([]string, 0, len(envSlice))
		for _, e := range envSlice {
			if !strings.HasPrefix(e, prefix) {
				filtered = append(filtered, e)
			}
		}
		// 新しい値を追加
		envSlice = filtered
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}
	
	fmt.Fprintf(os.Stderr, "%d個の環境変数を設定して指定されたコマンドを実行します: %s\n", count, strings.Join(cmdArgs, " "))
	
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = envSlice
	
	return cmd.Run()
}

func (c *CLI) runDump() error {
	data, err := file.ReadVaultedFile(c.vaultedFile)
	if err != nil {
		return fmt.Errorf(".env.vaultedファイルの読み込みに失敗しました: %w", err)
	}

	var password string
	if c.passwordStdin {
		password, err = utils.GetPasswordFromStdin()
	} else {
		password, err = utils.GetPasswordInteractive("復号化用パスワードを入力してください: ")
	}
	if err != nil {
		return err
	}

	decryptedData, err := crypto.Decrypt(data, password)
	if err != nil {
		return fmt.Errorf("復号化に失敗しました: %w", err)
	}

	fmt.Print(string(decryptedData))
	return nil
}

func (c *CLI) selectEnvironmentVariables(envVars []tui.EnvVar) ([]tui.EnvVar, error) {
	// デフォルトではBubbleteaを使用
	return tui.EnvVarSelection(envVars, tui.BubbleteaTUI)
}