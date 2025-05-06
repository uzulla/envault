package tui

// 将来的に別のTUIライブラリを追加する場合に備えて
// ProviderType は使用するTUIプロバイダの種類を表す列挙型です
type ProviderType int

const (
	// BubbleteaProvider は Bubbletea ライブラリを使用するプロバイダを表します
	BubbleteaTUI ProviderType = iota
	// 将来的に他のプロバイダを追加できます
)

// NewSelectionProvider は指定されたプロバイダタイプに基づいてSelectionProviderを作成します
func NewSelectionProvider(providerType ProviderType) SelectionProvider {
	// 現在はBubbleteaのみサポート
	return NewBubbleteaProvider()
}

// EnvVarSelection は環境変数選択UIを実行するためのヘルパー関数です
func EnvVarSelection(envVars []EnvVar, providerType ProviderType) ([]EnvVar, error) {
	provider := NewSelectionProvider(providerType)
	return provider.RunSelection(envVars)
}