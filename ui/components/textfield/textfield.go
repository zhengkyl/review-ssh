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
	tabBorder    = lipgloss.RoundedBorder()
	inputStyle   = lipgloss.NewStyle().Border(tabBorder, true) //.BorderBottom(true)
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
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

func (m *Model) Height() int {
	return m.props.Height
}
func (m *Model) Width() int {
	return m.props.Width
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	if m.focused && !m.inner.Focused() {
		return m, m.inner.Focus()
	}

	var cmd tea.Cmd
	m.inner, cmd = m.inner.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.props.Global.KeyMap.Back) {
			m.Blur()
		}

	}
	return m, cmd
}

func (m *Model) View() string {
	return inputStyle.Render(m.inner.View())
}

// textinput model
func (m *Model) Value() string {
	return m.inner.Value()
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
func (m *Model) BackgroundStyle(s lipgloss.Style) {
	m.inner.BackgroundStyle = s
}
func (m *Model) PlaceholderStyle(s lipgloss.Style) {
	m.inner.PlaceholderStyle = s
}
func (m *Model) CursorStyle(s lipgloss.Style) {
	m.inner.CursorStyle = s
}
