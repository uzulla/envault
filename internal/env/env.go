package env

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func ParseEnvContent(content []byte) (map[string]string, error) {
	envVars := make(map[string]string)
	scanner := bufio.NewScanner(bytes.NewReader(content))
	
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("行 %d: 無効な形式です。'KEY=VALUE'の形式が必要です", lineNum)
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
		   (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}
		
		envVars[key] = value
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ファイルの読み込み中にエラーが発生しました: %w", err)
	}
	
	return envVars, nil
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

func GenerateUnsetScript(envVars map[string]string) string {
	var script strings.Builder
	script.WriteString("#!/bin/bash\n\n")
	
	for key := range envVars {
		script.WriteString(fmt.Sprintf("unset %s\n", key))
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
