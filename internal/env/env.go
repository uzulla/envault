package env

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/uzulla/envault/internal/tui"
)

// .envファイルの内容をパースして環境変数のマップを返します
func ParseEnvContent(data []byte) (map[string]string, error) {
	envVars := make(map[string]string)
	
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// 空行やコメント行をスキップ
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// = がない行を不正として無視
		if !strings.Contains(line, "=") {
			continue
		}
		
		// Key=Value の形式を解析
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// クォーテーションを削除
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}
		
		envVars[key] = value
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	return envVars, nil
}

// .envファイルの内容をパースして環境変数リストを返します（コメント付き）
func ParseEnvContentWithComments(data []byte) ([]tui.EnvVar, error) {
	var result []tui.EnvVar
	var orderedKeys []string // 環境変数の出現順を保持
	var envVarsWithComments = make(map[string]tui.EnvVar) // キーと環境変数オブジェクトのマップ
	
	envVarsMap, err := ParseEnvContent(data)
	if err != nil {
		return nil, err
	}
	
	// コメントを再度スキャンして環境変数に関連付ける
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var lastComment string
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// コメント行を保存（複数行のコメントを連結）
		if strings.HasPrefix(line, "#") {
			comment := strings.TrimPrefix(line, "#")
			comment = strings.TrimSpace(comment)
			
			// 複数行のコメントを連結
			if lastComment != "" {
				lastComment += " " + comment
			} else {
				lastComment = comment
			}
			continue
		}
		
		// Key=Value の行を処理
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			
			// マップから値を取得し、出現順を保持
			if value, exists := envVarsMap[key]; exists {
				// 同じキーが複数回出現した場合は最後の値を使用
				// ただし順序リストには一度だけ追加
				alreadyAdded := false
				for _, existingKey := range orderedKeys {
					if existingKey == key {
						alreadyAdded = true
						break
					}
				}
				
				if !alreadyAdded {
					orderedKeys = append(orderedKeys, key)
				}
				
				// 環境変数オブジェクトを作成しマップに保存
				envVarsWithComments[key] = tui.EnvVar{
					Key:     key,
					Value:   value,
					Comment: lastComment,
					Enabled: true,
				}
				
				lastComment = "" // コメントをリセット
			}
		} else {
			lastComment = "" // Key=Value の形式でない行ではコメントをリセット
		}
	}
	
	// 出現順に基づいて結果を構築
	for _, key := range orderedKeys {
		if envVar, exists := envVarsWithComments[key]; exists {
			result = append(result, envVar)
		}
	}
	
	return result, nil
}

// エクスポート用のスクリプトを生成します
func GenerateExportScript(envVars map[string]string) string {
	var script strings.Builder
	script.WriteString("#!/bin/bash\n\n")
	
	for key, value := range envVars {
		if strings.ContainsAny(value, " \t\n\r\"'`$&|;<>(){}[]\\") {
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
		
		if strings.ContainsAny(ev.Value, " \t\n\r\"'`$&|;<>(){}[]\\") {
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

// 有効な環境変数の数をカウントします
func CountEnabledEnvVars(envVars []tui.EnvVar) int {
	count := 0
	for _, ev := range envVars {
		if ev.Enabled {
			count++
		}
	}
	return count
}

// アンセット用のスクリプトを生成します
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
		if !ev.Enabled {
			continue // 無効な環境変数はスキップ
		}
		
		script.WriteString(fmt.Sprintf("unset %s\n", ev.Key))
	}
	
	return script.String()
}

// シェルスクリプト内でのエスケープ処理を行います
func escapeValue(value string) string {
	return strings.ReplaceAll(value, "\"", "\\\"")
}