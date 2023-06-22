package skeleton

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
)

const SKELETON_FPS = 10
const NUM_FRAMES = 20

// Used to make sure only one skeleton sends ticks, which all skeletons use
var (
	lastID int
	idMtx  sync.Mutex
)

func nextID() int {
	idMtx.Lock()
	defer idMtx.Unlock()
	lastID++
	return lastID
}

type Model struct {
	props common.Props
	id    int
	frame int
}

func New(props common.Props) *Model {
	m := &Model{
		props: props,
		id:    nextID(),
	}

	return m
}

func (m *Model) SetSize(width, height int) {
	m.props.Width = width
	m.props.Height = height
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	switch msg.(type) {
	case TickMsg:
		m.frame = (m.frame + 1) % NUM_FRAMES

		return m, m.tick()
	default:
		return m, nil
	}
}

type TickMsg struct {
	Time time.Time
	ID   int
}

func (m *Model) View() string {
	rgb := loopInt(40, 100, float64(m.frame)/NUM_FRAMES)

	base := lipgloss.NewStyle().Background(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", rgb, rgb, rgb)))

	line := strings.Repeat(" ", m.props.Width) + "\n"
	view := strings.Repeat(line, m.props.Height)

	return base.Render(view[:len(view)-1])
}

func loopInt(min int, max int, frac float64) int {
	return int(math.Abs(0.5-frac)*float64(max-min)) + min
}

func (m *Model) Tick() tea.Msg {
	if m.id != 1 {
		return nil
	}

	return TickMsg{
		Time: time.Now(),
	}
}

func (m *Model) tick() tea.Cmd {
	if m.id != 1 {
		return nil
	}

	return tea.Tick(time.Second/SKELETON_FPS, func(t time.Time) tea.Msg {
		return TickMsg{
			Time: t,
		}
	})
}
