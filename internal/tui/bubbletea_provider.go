package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// BubbleteaProvider はBubbleteaライブラリを使用したSelectionProviderの実装です
type BubbleteaProvider struct{}

// 新しいBubbleteaProviderを作成します
func NewBubbleteaProvider() *BubbleteaProvider {
	return &BubbleteaProvider{}
}

// キーマッピングを定義
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Toggle key.Binding
	Help   key.Binding
	Quit   key.Binding
	Select key.Binding
}

// 使用可能なキーのキーマップを設定
var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "上へ"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "下へ"),
	),
	Toggle: key.NewBinding(
		key.WithKeys(" ", "space"),
		key.WithHelp("space", "選択/解除"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "ヘルプ"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc", "キャンセル"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "確定"),
	),
}

// ShortHelp は簡易ヘルプテキストを返します
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Toggle, k.Select, k.Quit}
}

// FullHelp はすべてのヘルプテキストを返します
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Toggle},
		{k.Select, k.Help, k.Quit},
	}
}

// モデルはBubbleteaアプリケーションの状態を表します
type model struct {
	envVars        []EnvVar
	cursor         int
	selected       map[string]bool
	viewport       viewport.Model
	help           help.Model
	keys           keyMap
	quitting       bool
	width          int
	height         int
	selectedResult []EnvVar
	showFullHelp   bool
}

// 注: スタイル設定は使用していません

// 初期化関数
func initialModel(envVars []EnvVar) model {
	helpModel := help.New()
	selected := make(map[string]bool)

	// デフォルトですべての環境変数を有効にする
	for i, ev := range envVars {
		envVars[i].Enabled = true
		selected[ev.Key] = true // 初期状態ではすべて選択済み
	}

	// デフォルトサイズを設定
	width := 80
	height := 30

	// 項目数と画面サイズに応じてviewportの高さを決定
	viewportHeight := len(envVars) + 2
	if viewportHeight > height-4 {
		viewportHeight = height - 4
	}
	if viewportHeight < 1 {
		viewportHeight = 1
	}

	vp := viewport.New(width, viewportHeight)
	vp.SetContent(renderItems(envVars, selected, 0))

	return model{
		envVars:      envVars,
		cursor:       0,
		selected:     selected,
		viewport:     vp,
		help:         helpModel,
		keys:         keys,
		width:        width,
		height:       height,
		showFullHelp: false,
	}
}

// カーソルマークを返す
func cursorMark(i, cursor int) string {
	if i == cursor {
		return ">"
	}
	return " "
}

// チェックボックスマークを返す
func checkMark(selected bool) string {
	if selected {
		// 選択されている場合は明確なチェックマーク
		return "[✓]"
	}
	// 選択されていない場合は空のボックス
	return "[ ]"
}

// アイテムリストをレンダリング
func renderItems(envVars []EnvVar, selected map[string]bool, cursor int) string {
	var b strings.Builder

	b.WriteString("適用する環境変数を選択してください\n\n")

	// すべての項目を完全に固定フォーマットでレンダリング（スタイル無し）
	for i, ev := range envVars {
		isSelected := selected[ev.Key]
		
		// 固定の左マージン（多くの行で同じマージンを使用）
		b.WriteString("   ")
		
		// カーソル表示（スタイル無し）
		if i == cursor {
			b.WriteString("> ")
		} else {
			b.WriteString("  ")
		}
		
		// チェックボックス - 確実に現在の選択状態を反映
		check := checkMark(isSelected)
		// 選択状態がわかりやすいよう、チェックボックスとその前後にスペースを入れる
		b.WriteString(check + " ")
		
		// 環境変数名（スタイル無し）
		b.WriteString(ev.Key)
		
		// コメント（スタイル無し）
		if ev.Comment != "" {
			b.WriteString(" - " + ev.Comment)
		}
		
		b.WriteString("\n")
	}

	return b.String()
}

// Init は初期コマンドを返します
func (m model) Init() tea.Cmd {
	return nil
}

// Update はモデル状態を更新します
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
				// 選択項目が表示領域外に移動した場合、スクロール位置を調整
				if m.cursor < m.viewport.YOffset {
					m.viewport.SetYOffset(m.cursor)
				}
			}

		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.envVars)-1 {
				m.cursor++
				// 選択項目が表示領域外に移動した場合、スクロール位置を調整
				if m.cursor >= m.viewport.YOffset+m.viewport.Height {
					m.viewport.SetYOffset(m.cursor - m.viewport.Height + 1)
				}
			}

		case key.Matches(msg, m.keys.Help):
			// ヘルプ表示の切り替え
			m.showFullHelp = !m.showFullHelp

		case key.Matches(msg, m.keys.Toggle):
			if len(m.envVars) == 0 {
				break
			}
			// カーソル位置の環境変数のキーを取得
			key := m.envVars[m.cursor].Key
			// 現在の状態を取得し、確実に反転する
			currentState, ok := m.selected[key]
			if !ok {
				// マップに存在しない場合は初期化
				currentState = true // デフォルトで有効
			}
			newState := !currentState
			// マップとEnvVar構造体の両方を更新
			m.selected[key] = newState
			m.envVars[m.cursor].Enabled = newState
			
			// 更新をUIに反映させるために明示的にコンテンツを再設定
			m.viewport.SetContent(renderItems(m.envVars, m.selected, m.cursor))

		case key.Matches(msg, m.keys.Select):
			// 選択された環境変数を結果として設定
			var selectedEnvVars []EnvVar
			for i, ev := range m.envVars {
				// 強制的にマップの値を使用して設定
				isSelected := m.selected[ev.Key]
				m.envVars[i].Enabled = isSelected
				selectedEnvVars = append(selectedEnvVars, m.envVars[i])
			}
			m.selectedResult = selectedEnvVars
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		
		// ヘルプ表示のためのスペースを確保し、最小値を1に制限
		viewportHeight := msg.Height - 4
		if viewportHeight < 1 {
			viewportHeight = 1
		}
		m.viewport.Height = viewportHeight
		m.help.Width = msg.Width
	}

	m.viewport.SetContent(renderItems(m.envVars, m.selected, m.cursor))
	return m, nil
}

// View はUIをレンダリングします
func (m model) View() string {
	if m.quitting {
		return ""
	}

	var helpView string
	if m.showFullHelp {
		helpView = m.help.FullHelpView(m.keys.FullHelp())
	} else {
		helpView = m.help.View(m.keys)
	}

	return fmt.Sprintf("%s\n\n%s", m.viewport.View(), helpView)
}

// RunSelection は環境変数選択UIを実行し、選択された環境変数を返します
func (b *BubbleteaProvider) RunSelection(envVars []EnvVar) ([]EnvVar, error) {
	p := tea.NewProgram(initialModel(envVars), tea.WithAltScreen())
	
	m, err := p.Run()
	if err != nil {
		return nil, err
	}
	
	if m, ok := m.(model); ok {
		if m.quitting || m.selectedResult == nil || len(m.selectedResult) == 0 {
			// キャンセルされた場合、またはresultが空の場合は元の環境変数をすべて有効にして返す
			for i := range envVars {
				envVars[i].Enabled = true
			}
			return envVars, nil
		}
		return m.selectedResult, nil
	}
	
	return envVars, nil
}