package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/uzulla/envault/internal/crypto"
	"github.com/uzulla/envault/internal/env"
	"github.com/uzulla/envault/internal/file"
	"github.com/uzulla/envault/pkg/utils"
)

const (
	Version = "0.1.0"
)

var (
	ErrInvalidCommand = errors.New("無効なコマンドです")
)

type CLI struct {
	passwordStdin bool
	vaultedFile   string
	newShell      bool
}

func NewCLI() *CLI {
	return &CLI{}
}

func (c *CLI) Run(args []string) error {
	if len(args) < 1 {
		c.printUsage()
		return nil
	}

	globalFlags := flag.NewFlagSet("envault", flag.ExitOnError)
	globalFlags.BoolVar(&c.passwordStdin, "password-stdin", false, "stdinからパスワードを読み込む")
	globalFlags.StringVar(&c.vaultedFile, "file", "", "使用する.env.vaultedファイルのパス")
	
	var outputScriptOnly bool
	globalFlags.BoolVar(&outputScriptOnly, "output-script-only", false, "スクリプトのみを出力（情報メッセージなし）")
	globalFlags.BoolVar(&c.newShell, "new-shell", false, "新しいbashセッションを起動して環境変数を設定")

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
			return c.runExport(outputScriptOnly, cmdArgs)
		} else {
			if err := globalFlags.Parse(args[1:]); err != nil {
				return err
			}
			return c.runExport(outputScriptOnly, nil)
		}
	case "unset":
		if err := globalFlags.Parse(args[1:]); err != nil {
			return err
		}
		return c.runUnset(outputScriptOnly, nil)
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

func (c *CLI) runExport(outputScriptOnly bool, cmdArgs []string) error {
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

	envVars, err := env.ParseEnvContent(decryptedData)
	if err != nil {
		return fmt.Errorf("環境変数の解析に失敗しました: %w", err)
	}

	if c.newShell {
		for k, v := range envVars {
			os.Setenv(k, v)
		}
		
		fmt.Fprintf(os.Stderr, "%d個の環境変数を設定して新しいbashセッションを起動します\n", len(envVars))
		
		cmd := exec.Command("/bin/bash")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		
		return cmd.Run()
	} else if len(cmdArgs) > 0 {
		for k, v := range envVars {
			os.Setenv(k, v)
		}
		
		fmt.Fprintf(os.Stderr, "%d個の環境変数を設定して指定されたコマンドを実行します: %s\n", len(envVars), strings.Join(cmdArgs, " "))
		
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		
		return cmd.Run()
	} else if outputScriptOnly {
		script := env.GenerateExportScript(envVars)
		fmt.Print(script)
	} else {
		script := env.GenerateExportScript(envVars)
		if err := utils.ExecuteScript(script); err != nil {
			return fmt.Errorf("環境変数のエクスポートに失敗しました: %w", err)
		}
		fmt.Fprintf(os.Stderr, "%d個の環境変数をエクスポートしました\n", len(envVars))
	}

	return nil
}

func (c *CLI) runUnset(outputScriptOnly bool, cmdArgs []string) error {
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

	envVars, err := env.ParseEnvContent(decryptedData)
	if err != nil {
		return fmt.Errorf("環境変数の解析に失敗しました: %w", err)
	}

	script := env.GenerateUnsetScript(envVars)

	if outputScriptOnly {
		fmt.Print(script)
	} else {
		if err := utils.ExecuteScript(script); err != nil {
			return fmt.Errorf("環境変数のアンセットに失敗しました: %w", err)
		}
		fmt.Fprintf(os.Stderr, "%d個の環境変数をアンセットしました\n", len(envVars))
	}

	return nil
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
  
  echo "password" | envault export --password-stdin  stdinからパスワードを読み込んでエクスポート
  
  # 環境変数をアンセットする方法:
  envault unset                       アンセットスクリプトのパスを表示
  eval $(envault unset --output-script-only)  環境変数を直接アンセット
  source <(envault unset --output-script-only)  環境変数を直接アンセット（別の方法）
  
  # 暗号化されたファイルの内容を確認する方法:
  envault dump                        .env.vaultedファイルの内容を表示
  envault dump > decrypted.env        復号化した内容をファイルに保存`)
}
