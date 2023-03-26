package overlay

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type OverlayModel struct {
	common       common.Common
	parentWidth  int
	parentHeight int
	child        tea.Model
}

func New(common common.Common, parentWidth int, parentHeight int, child tea.Model) *OverlayModel {
	return &OverlayModel{
		common:       common,
		parentWidth:  parentWidth,
		parentHeight: parentHeight,
		child:        child,
	}
}

func (m *OverlayModel) Update(msg tea.Msg) (*OverlayModel, tea.Cmd) {

	var cmd tea.Cmd

	// TODO should we accept dirty non-pointer children?
	m.child, cmd = m.child.Update(msg)
	return m, cmd
}

// func (m *OverlayModel) view() string {

// 	// todo padding and margin + etc

// 	return m.child.View()
// }

func (m *OverlayModel) RenderOverlay(parentView string) string {

	// todo padding and margin + etc

	strings.Split(parentView, "\n")

	return m.child.View()
}
