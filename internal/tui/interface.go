package tui

// EnvVar は環境変数のキー、値、およびコメントを表す構造体です
type EnvVar struct {
	Key     string
	Value   string
	Comment string
	Enabled bool
}

// SelectionProvider はTUIセレクションの実装を提供するインターフェースです
// これにより、将来的にTUIライブラリを交換できるようになります
type SelectionProvider interface {
	// RunSelection は環境変数選択UIを実行し、選択された環境変数を返します
	RunSelection(envVars []EnvVar) ([]EnvVar, error)
}