package skeleton

import (
	"fmt"
	"math"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
)

const skeletonFPS = 10

const animationFrames = 20

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

type SkeletonModel struct {
	common common.Common
	id     int
	frame  int
}

func New(common common.Common) *SkeletonModel {
	m := &SkeletonModel{
		common: common,
		id:     nextID(),
	}

	return m
}

func (m SkeletonModel) Update(msg tea.Msg) (SkeletonModel, tea.Cmd) {
	switch msg.(type) {
	case TickMsg:
		m.frame = (m.frame + 1) % animationFrames

		return m, m.tick()
	default:
		return m, nil
	}
}

type TickMsg struct {
	Time time.Time
	ID   int
}

func (m SkeletonModel) View() string {
	rgb := loopInt(40, 100, float64(m.frame)/animationFrames)

	base := lipgloss.NewStyle().Background(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", rgb, rgb, rgb)))

	view := ""

	for y := 0; y < m.common.Height; y++ {
		for x := 0; x < m.common.Width; x++ {
			view += ("　")
			// view += fmt.Sprint(m.frame)
		}
		view += "\n"
	}

	return base.Render(view)
}

func loopInt(min int, max int, frac float64) int {
	return int(math.Abs(0.5-frac)*float64(max-min)) + min
}

func (m SkeletonModel) Tick() tea.Msg {
	if m.id != 1 {
		return nil
	}

	return TickMsg{
		Time: time.Now(),
	}
}

func (m SkeletonModel) tick() tea.Cmd {
	if m.id != 1 {
		return nil
	}

	return tea.Tick(time.Second/skeletonFPS, func(t time.Time) tea.Msg {
		return TickMsg{
			Time: t,
		}
	})
}
