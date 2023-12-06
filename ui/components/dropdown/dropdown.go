package dropdown

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/util"
)

var (
	tabBorder      = lipgloss.RoundedBorder()
	unfocusedStyle = lipgloss.NewStyle().Border(tabBorder, true)
	dividerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F25D94"))
	closedStyle    = lipgloss.NewStyle().Border(tabBorder, true).BorderForeground(lipgloss.Color("#F25D94"))
	openStyle      = lipgloss.NewStyle().Border(tabBorder, true).BorderBottom(false).BorderForeground(lipgloss.Color("#F25D94"))
	itemStyle      = lipgloss.NewStyle().Border(tabBorder, false, true).BorderForeground(lipgloss.Color("#F25D94"))
	lastStyle      = lipgloss.NewStyle().Border(tabBorder, false, true, true).BorderForeground(lipgloss.Color("#F25D94"))
	activeStyle    = lipgloss.NewStyle().Background(lipgloss.Color("#F25D94")).Padding(0, 1)
	normalStyle    = lipgloss.NewStyle().Padding(0, 1)
)

type onChange func(value string) tea.Cmd

type Option struct {
	Text  string
	Value string
}

type Model struct {
	props    common.Props
	focused  bool
	OnChange onChange
	open     bool
	noneText string
	options  []Option
	Selected int // -1 if none, else index into options
	active   int // not yet selected, but hovering index
}

func New(p common.Props, noneText string, options []Option) *Model {
	return &Model{
		props:    p,
		noneText: noneText,
		options:  options,
		focused:  false,
		open:     false,
		Selected: -1,
		active:   -1,
	}
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
	m.open = false
}

func (m *Model) SetSize(width, height int) {
	m.props.Width = width
	m.props.Height = height
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *common.KeyEvent:
		switch {
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Back):
			if m.open {
				msg.Handled = true
				m.open = false
			}
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Down):
			msg.Handled = true
			m.active = util.Min(m.active+1, len(m.options)-1)
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Up):
			msg.Handled = true
			m.active = util.Max(m.active-1, 0)
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Select):
			msg.Handled = true
			if !m.open {
				m.open = true
				if m.Selected == -1 {
					m.active = 0
				} else {
					m.active = m.Selected
				}

			} else {
				m.Selected = m.active
				m.open = false
				return m, m.OnChange(m.options[m.Selected].Value)
			}
		}
	}
	return m, nil
}

func (m *Model) View() string {
	hf := unfocusedStyle.GetHorizontalFrameSize()
	itemHf := itemStyle.GetHorizontalFrameSize()

	itemWidth := m.props.Width - hf - itemHf

	var selected string
	if m.Selected == -1 {
		selected = m.noneText
	} else {
		selected = m.options[m.Selected].Text
	}
	selected = util.TruncOrPadASCII(selected, itemWidth-2) + " ▼"
	selected = " " + selected + " "

	if !m.open {
		if m.focused {
			return closedStyle.Render(selected)
		} else {
			return unfocusedStyle.Render(selected)
		}
	}

	sb := strings.Builder{}
	sb.WriteString(openStyle.Render(selected))

	sb.WriteString("\n" + dividerStyle.Render("├"+strings.Repeat("─", itemWidth+2)+"┤") + "\n")

	for i, option := range m.options {
		text := util.TruncOrPadASCII(option.Text, itemWidth)
		if i == m.active {
			text = activeStyle.Render(text)
		} else {
			text = normalStyle.Render(text)
		}

		if i == len(m.options)-1 {
			sb.WriteString(lastStyle.Render(text))
		} else {
			sb.WriteString(itemStyle.Render(text))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
