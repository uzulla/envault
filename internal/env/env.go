package env

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
)

func ParseEnvContent(content []byte) (map[string]string, error) {
	envVars, err := godotenv.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf(".envファイルの解析に失敗しました: %w", err)
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
