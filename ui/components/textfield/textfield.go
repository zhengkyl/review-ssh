package textfield

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/reflow/ansi"
	"github.com/zhengkyl/review-ssh/ui/common"
)

var (
	tabBorder  = lipgloss.RoundedBorder()
	inputStyle = lipgloss.NewStyle().Border(tabBorder, true)
)

type Model struct {
	common      common.Common
	inner       textinput.Model
	placeholder string
}

func New(common common.Common) *Model {
	inner := textinput.New()
	// inner.Width = common.Width - 8

	return &Model{common, inner, ""}
}

func (m *Model) Focus() tea.Cmd {
	return m.inner.Focus()
}

func (m *Model) Blur() {
	m.inner.Blur()
}

func (m *Model) SetSize(w, h int) {
	// TODO figure out how to get this padding/margin

	m.common.Width = w - 8
	m.common.Height = h

	m.inner.Width = w - 8

	if m.placeholder != "" {
		m.Placeholder(m.placeholder)
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.inner, cmd = m.inner.Update(msg)

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
