package cli

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/uzulla/envault/internal/crypto"
	"github.com/uzulla/envault/internal/env"
	"github.com/uzulla/envault/internal/file"
	"github.com/uzulla/envault/internal/tui"
	"github.com/uzulla/envault/pkg/utils"
)

const (
	Version = "0.1.0"
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
	passwordStdin bool
	vaultedFile   string
	newShell      bool
	selectVars    bool // 環境変数を選択するオプション
}

func NewCLI() *CLI {
	return &CLI{}
}

func (c *CLI) Run(args []string) error {
	if len(args) < 1 {
		c.printUsage()
		return nil
	}

	// フラグ定義 - パース時の問題を避けるためにContinueOnErrorを使用
	globalFlags := flag.NewFlagSet("envault", flag.ContinueOnError)
	globalFlags.SetOutput(ioutil.Discard) // パースエラーは呼び出し側で扱う

	// フラグの設定
	globalFlags.BoolVar(&c.passwordStdin, "password-stdin", false, "stdinからパスワードを読み込む")
	globalFlags.StringVar(&c.vaultedFile, "file", "", "使用する.env.vaultedファイルのパス") 
	
	var outputScriptOnly bool
	globalFlags.BoolVar(&outputScriptOnly, "output-script-only", false, "スクリプトのみを出力（情報メッセージなし）")
	globalFlags.BoolVar(&c.newShell, "new-shell", false, "新しいbashセッションを起動して環境変数を設定")
	
	// selectフラグを明示的に定義
	globalFlags.BoolVar(&c.selectVars, "select", false, "適用する環境変数をTUIで選択する")

	command := args[0]

	switch command {
	case "help", "-h", "--help":
		c.printUsage()
		return nil
	case "version", "-v", "--version":
		fmt.Printf("envault version %s\n", Version)
		return nil
	case "export":
		var cmdArgs []string
		dashDashIndex := -1
		for i, arg := range args {
			if arg == "--" && i > 0 {
				dashDashIndex = i
				break
			}
		}
		
		if dashDashIndex != -1 {
			cmdArgs = args[dashDashIndex+1:]
			if err := globalFlags.Parse(args[1:dashDashIndex]); err != nil {
				return err
			}
			return c.runWithVaultedFile(ExportMode, outputScriptOnly, cmdArgs)
		} else {
			if err := globalFlags.Parse(args[1:]); err != nil {
				return err
			}
			return c.runWithVaultedFile(ExportMode, outputScriptOnly, nil)
		}
	case "unset":
		if err := globalFlags.Parse(args[1:]); err != nil {
			return err
		}
		return c.runWithVaultedFile(UnsetMode, outputScriptOnly, nil)
	case "dump":
		if err := globalFlags.Parse(args[1:]); err != nil {
			return err
		}
		return c.runDump()
	default:
		if err := globalFlags.Parse(args[1:]); err != nil {
			return err
		}
		return c.runEncrypt(command)
	}
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

func (c *CLI) printUsage() {
	fmt.Println(`使用方法:
  envault [オプション] <.envファイル>  .envファイルを暗号化して.env.vaultedファイルを作成
  envault export [オプション]          .env.vaultedファイルから環境変数をエクスポート
  envault export [オプション] -- <コマンド> [引数...]  環境変数を設定して指定したコマンドを実行
  envault unset [オプション]           .env.vaultedファイルに記載された環境変数をアンセット
  envault dump [オプション]            .env.vaultedファイルを復号化して内容を表示
  envault help                        ヘルプを表示
  envault version                     バージョン情報を表示

オプション:
  --password-stdin                    stdinからパスワードを読み込む
  --file <ファイルパス>               使用する.env.vaultedファイルのパス（デフォルト: .env.vaulted）
  --output-script-only                スクリプトのみを出力（情報メッセージなし）
  --new-shell                         新しいbashセッションを起動して環境変数を設定
  --select                            適用する環境変数をTUIで選択する

例:
  envault .env                        .envファイルを暗号化
  
  # 環境変数をエクスポートする方法:
  envault export                      エクスポートスクリプトのパスを表示
  eval $(envault export --output-script-only)  環境変数を直接エクスポート
  source <(envault export --output-script-only)  環境変数を直接エクスポート（別の方法）
  
  # 新しい方法で環境変数を使用する:
  envault export --new-shell          新しいbashセッションを起動して環境変数を設定
  envault export -- node app.js       環境変数を設定してnodeコマンドを実行
  envault export -- docker-compose up 環境変数を設定してdocker-composeを実行
  
  # 環境変数を選択的に適用する方法:
  envault export --select             TUIで環境変数を選択してからエクスポート
  envault export --select --new-shell TUIで環境変数を選択してから新しいbashセッションを起動
  envault export --select -- npm start TUIで環境変数を選択してからコマンドを実行
  
  echo "password" | envault export --password-stdin  stdinからパスワードを読み込んでエクスポート
  
  # 環境変数をアンセットする方法:
  envault unset                       アンセットスクリプトのパスを表示
  eval $(envault unset --output-script-only)  環境変数を直接アンセット
  source <(envault unset --output-script-only)  環境変数を直接アンセット（別の方法）
  
  # 暗号化されたファイルの内容を確認する方法:
  envault dump                        .env.vaultedファイルの内容を表示
  envault dump > decrypted.env        復号化した内容をファイルに保存`)
}