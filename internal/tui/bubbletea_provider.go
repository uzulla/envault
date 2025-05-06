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
	h := help.New()
	selected := make(map[string]bool)

	// デフォルトですべての環境変数を有効にする
	for i, ev := range envVars {
		envVars[i].Enabled = true
		selected[ev.Key] = true
	}

	vp := viewport.New(80, len(envVars)+2)
	vp.SetContent(renderItems(envVars, selected, 0))

	return model{
		envVars:  envVars,
		cursor:   0,
		selected: selected,
		viewport: vp,
		help:     h,
		keys:     keys,
		width:    80,
		height:   30,
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
			}

		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.envVars)-1 {
				m.cursor++
			}

		case key.Matches(msg, m.keys.Toggle):
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
		m.viewport.Height = msg.Height - 4 // ヘルプ表示のためのスペースを確保
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

	return fmt.Sprintf("%s\n\n%s", m.viewport.View(), m.help.View(m.keys))
}

// RunSelection は環境変数選択UIを実行し、選択された環境変数を返します
func (b *BubbleteaProvider) RunSelection(envVars []EnvVar) ([]EnvVar, error) {
	p := tea.NewProgram(initialModel(envVars), tea.WithAltScreen())
	
	m, err := p.Run()
	if err != nil {
		return nil, err
	}
	
	if m, ok := m.(model); ok {
		if m.quitting && m.selectedResult == nil {
			// キャンセルされた場合は元の環境変数をすべて有効にして返す
			for i := range envVars {
				envVars[i].Enabled = true
			}
			return envVars, nil
		}
		return m.selectedResult, nil
	}
	
	return envVars, nil
}