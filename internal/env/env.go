package env

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/uzulla/envault/internal/tui"
)

// 環境変数をキーと値のマップとして解析します
func ParseEnvContent(content []byte) (map[string]string, error) {
	envVars, err := godotenv.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf(".envファイルの解析に失敗しました: %w", err)
	}
	
	return envVars, nil
}

// コメントを含む環境変数をEnvVar構造体のスライスとして解析します
func ParseEnvContentWithComments(content []byte) ([]tui.EnvVar, error) {
	// 基本的な環境変数の解析
	envVarsMap, err := godotenv.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf(".envファイルの解析に失敗しました: %w", err)
	}
	
	// 先行コメントを保持するマップ
	commentMap := make(map[string]string)
	
	// 行ごとに解析してコメントを抽出
	scanner := bufio.NewScanner(bytes.NewReader(content))
	var lastComment string
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// コメント行の場合
		if strings.HasPrefix(line, "#") {
			lastComment = strings.TrimSpace(strings.TrimPrefix(line, "#"))
			continue
		}
		
		// 空行の場合はコメントをリセット
		if line == "" {
			lastComment = ""
			continue
		}
		
		// 環境変数の行の場合
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			
			// コメントがあればマップに保存
			if lastComment != "" {
				commentMap[key] = lastComment
				lastComment = ""
			}
		}
	}
	
	// 環境変数とコメントを結合
	var result []tui.EnvVar
	for key, value := range envVarsMap {
		result = append(result, tui.EnvVar{
			Key:     key,
			Value:   value,
			Comment: commentMap[key],
			Enabled: true, // デフォルトで有効
		})
	}
	
	return result, nil
}

func GenerateExportScript(envVars map[string]string) string {
	var script strings.Builder
	script.WriteString("#!/bin/bash\n\n")
	
	for key, value := range envVars {
		if strings.ContainsAny(value, " \t\n\r\"'`$&|;<>(){}[]") {
			script.WriteString(fmt.Sprintf("export %s=\"%s\"\n", key, escapeValue(value)))
		} else {
			script.WriteString(fmt.Sprintf("export %s=%s\n", key, value))
		}
	}
	
	return script.String()
}

// TUI選択後の環境変数リストからエクスポートスクリプトを生成します
func GenerateExportScriptFromEnvVarList(envVars []tui.EnvVar) string {
	var script strings.Builder
	script.WriteString("#!/bin/bash\n\n")
	
	for _, ev := range envVars {
		if !ev.Enabled {
			continue // 無効な環境変数はスキップ
		}
		
		if strings.ContainsAny(ev.Value, " \t\n\r\"'`$&|;<>(){}[]") {
			script.WriteString(fmt.Sprintf("export %s=\"%s\"\n", ev.Key, escapeValue(ev.Value)))
		} else {
			script.WriteString(fmt.Sprintf("export %s=%s\n", ev.Key, ev.Value))
		}
	}
	
	return script.String()
}

// TUI選択後の環境変数リストから有効な環境変数のマップを生成します
func FilterEnabledEnvVars(envVars []tui.EnvVar) map[string]string {
	result := make(map[string]string)
	
	for _, ev := range envVars {
		if ev.Enabled {
			result[ev.Key] = ev.Value
		}
	}
	
	return result
}

// 有効な環境変数の数を数えます
func CountEnabledEnvVars(envVars []tui.EnvVar) int {
	count := 0
	for _, ev := range envVars {
		if ev.Enabled {
			count++
		}
	}
	return count
}

func GenerateUnsetScript(envVars map[string]string) string {
	var script strings.Builder
	script.WriteString("#!/bin/bash\n\n")
	
	for key := range envVars {
		script.WriteString(fmt.Sprintf("unset %s\n", key))
	}
	
	return script.String()
}

// TUI選択後の環境変数リストからアンセットスクリプトを生成します
func GenerateUnsetScriptFromEnvVarList(envVars []tui.EnvVar) string {
	var script strings.Builder
	script.WriteString("#!/bin/bash\n\n")
	
	for _, ev := range envVars {
		if ev.Enabled {
			script.WriteString(fmt.Sprintf("unset %s\n", ev.Key))
		}
	}
	
	return script.String()
}

func escapeValue(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "`", "\\`")
	value = strings.ReplaceAll(value, "$", "\\$")
	
	return value
}