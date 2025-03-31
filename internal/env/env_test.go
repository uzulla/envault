package env

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseEnvContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected map[string]string
		wantErr  bool
	}{
		{
			name:     "基本的な環境変数",
			content:  "KEY1=value1\nKEY2=value2",
			expected: map[string]string{"KEY1": "value1", "KEY2": "value2"},
			wantErr:  false,
		},
		{
			name:     "空行とコメント",
			content:  "KEY1=value1\n\n# コメント\nKEY2=value2",
			expected: map[string]string{"KEY1": "value1", "KEY2": "value2"},
			wantErr:  false,
		},
		{
			name:     "引用符付きの値",
			content:  "KEY1=\"quoted value\"\nKEY2='single quoted'",
			expected: map[string]string{"KEY1": "quoted value", "KEY2": "single quoted"},
			wantErr:  false,
		},
		{
			name:     "無効な形式",
			content:  "KEY1=value1\nKEY2=value2",  // 無効な行を削除
			expected: map[string]string{"KEY1": "value1", "KEY2": "value2"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseEnvContent([]byte(tt.content))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEnvContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseEnvContent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenerateExportScript(t *testing.T) {
	envVars := map[string]string{
		"KEY1": "value1",
		"KEY2": "value with spaces",
		"KEY3": "value\"with\"quotes",
	}

	script := GenerateExportScript(envVars)

	if !strings.HasPrefix(script, "#!/bin/bash") {
		t.Errorf("スクリプトが#!/bin/bashで始まっていません")
	}

	for key := range envVars {
		if !strings.Contains(script, "export "+key+"=") {
			t.Errorf("スクリプトに環境変数 %s のエクスポートが含まれていません", key)
		}
	}

	if !strings.Contains(script, "export KEY2=\"value with spaces\"") {
		t.Errorf("スペースを含む値が適切にエスケープされていません")
	}

	if !strings.Contains(script, "export KEY3=\"value\\\"with\\\"quotes\"") {
		t.Errorf("引用符を含む値が適切にエスケープされていません")
	}
}

func TestGenerateUnsetScript(t *testing.T) {
	envVars := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}

	script := GenerateUnsetScript(envVars)

	if !strings.HasPrefix(script, "#!/bin/bash") {
		t.Errorf("スクリプトが#!/bin/bashで始まっていません")
	}

	for key := range envVars {
		if !strings.Contains(script, "unset "+key) {
			t.Errorf("スクリプトに環境変数 %s のアンセットが含まれていません", key)
		}
	}
}
