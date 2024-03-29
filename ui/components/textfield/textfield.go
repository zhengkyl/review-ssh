package textfield

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/reflow/ansi"
	"github.com/zhengkyl/review-ssh/ui/common"
)

var (
	tabBorder         = lipgloss.RoundedBorder()
	inputStyle        = lipgloss.NewStyle().Border(tabBorder, true)
	focusedInputStyle = lipgloss.NewStyle().Border(tabBorder, true).BorderForeground(lipgloss.Color("#F25D94"))
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#F25D94"))
	// blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	// cursorStyle  = focusedStyle.Copy()
	noStyle = lipgloss.NewStyle()
)

type Model struct {
	props       common.Props
	inner       textinput.Model
	focused     bool
	placeholder string
}

func New(p common.Props) *Model {
	inner := textinput.New()

	m := &Model{p, inner, false, ""}

	m.SetSize(p.Width, p.Height)

	return m
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) Focus() {
	m.focused = true
	m.inner.PromptStyle = focusedStyle
	m.inner.TextStyle = focusedStyle
}

func (m *Model) Blur() {
	m.focused = false
	m.inner.Blur()
	m.inner.PromptStyle = noStyle
	m.inner.TextStyle = noStyle
}

func (m *Model) SetSize(w, h int) {
	m.props.Width = w
	m.props.Height = h

	// Left right border + padding + > indicator
	m.inner.Width = w - 5

	if m.placeholder != "" {
		m.Placeholder(m.placeholder)
	}
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	var cmds []tea.Cmd
	if m.focused && !m.inner.Focused() {
		cmds = append(cmds, m.inner.Focus())
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case *common.KeyEvent:
		if key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Back) {
			m.Blur()
			msg.Handled = true
		} else {
			prevValue := m.inner.Value()
			prevPos := m.inner.Position()
			m.inner, cmd = m.inner.Update(msg.KeyMsg)
			if m.inner.Value() != prevValue ||
				m.inner.Position() != prevPos {
				msg.Handled = true
			}
		}
	default:
		m.inner, cmd = m.inner.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if m.focused {
		return focusedInputStyle.Render(m.inner.View())
	} else {
		return inputStyle.Render(m.inner.View())
	}
}

// textinput model
func (m *Model) Value() string {
	return m.inner.Value()
}

func (m *Model) SetValue(s string) {
	m.inner.SetValue(s)
}

func (m *Model) Prompt(p string) {
	m.inner.Prompt = p
}
func (m *Model) Placeholder(p string) {
	m.placeholder = p

	m.inner.Placeholder = m.placeholder

	phWidth := ansi.PrintableRuneWidth(m.inner.Placeholder)
	phBytes := len(m.inner.Placeholder)

	for phWidth > m.inner.Width {
		r, b := utf8.DecodeLastRuneInString(m.inner.Placeholder)
		phWidth -= runewidth.RuneWidth(r)
		phBytes -= b
		m.inner.Placeholder = m.inner.Placeholder[:phBytes]
	}
	// extra space after placeholder necessary to maintain same width after editing
	m.inner.Placeholder = m.inner.Placeholder + strings.Repeat(" ", m.inner.Width-phWidth+1)
}

func (m *Model) EchoMode(e textinput.EchoMode) {
	m.inner.EchoMode = e
}
func (m *Model) EchoCharacter(e rune) {
	m.inner.EchoCharacter = e
}
func (m *Model) CharLimit(c int) {
	m.inner.CharLimit = c
}
func (m *Model) Cursor(c cursor.Model) {
	m.inner.Cursor = c
}
func (m *Model) PromptStyle(s lipgloss.Style) {
	m.inner.PromptStyle = s
}
func (m *Model) TextStyle(s lipgloss.Style) {
	m.inner.TextStyle = s
}
func (m *Model) PlaceholderStyle(s lipgloss.Style) {
	m.inner.PlaceholderStyle = s
}
