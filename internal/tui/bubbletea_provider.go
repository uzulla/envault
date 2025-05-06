package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
		key.WithKeys("space"),
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

// スタイル設定
var (
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	cursorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	normalItemStyle   = lipgloss.NewStyle()
	checkedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	uncheckedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	commentStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

// 初期化関数
func initialModel(envVars []EnvVar) model {
	helpModel := help.New()
	selected := make(map[string]bool)

	// デフォルトですべての環境変数を有効にする
	for i, ev := range envVars {
		envVars[i].Enabled = true
		selected[ev.Key] = true
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
		return cursorStyle.Render("> ")
	}
	return "  "
}

// チェックボックスマークを返す
func checkMark(selected bool) string {
	if selected {
		return checkedStyle.Render("[✓]")
	}
	return uncheckedStyle.Render("[ ]")
}

// アイテムリストをレンダリング
func renderItems(envVars []EnvVar, selected map[string]bool, cursor int) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("適用する環境変数を選択してください\n\n"))

	for i, ev := range envVars {
		isSelected := selected[ev.Key]
		cursorMrk := cursorMark(i, cursor)
		check := checkMark(isSelected)

		// 環境変数の表示
		keyValueText := fmt.Sprintf("%s %s %s", cursorMrk, check, ev.Key)
		
		// カーソル位置にあるアイテムをハイライト
		if i == cursor {
			keyValueText = selectedItemStyle.Render(keyValueText)
		} else {
			keyValueText = normalItemStyle.Render(keyValueText)
		}

		b.WriteString(keyValueText)

		// コメントがある場合は表示
		if ev.Comment != "" {
			b.WriteString(" " + commentStyle.Render("- "+ev.Comment))
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
			key := m.envVars[m.cursor].Key
			m.selected[key] = !m.selected[key]
			m.envVars[m.cursor].Enabled = m.selected[key]

		case key.Matches(msg, m.keys.Select):
			// 選択された環境変数を結果として設定
			var selectedEnvVars []EnvVar
			for i, ev := range m.envVars {
				m.envVars[i].Enabled = m.selected[ev.Key]
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