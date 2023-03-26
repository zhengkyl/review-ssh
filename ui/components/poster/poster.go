package poster

import (
	"fmt"
	"image"

	"image/color"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/skeleton"
	"golang.org/x/image/draw"

	// _ means imported for its initialization side-effect
	_ "image/jpeg"
	_ "image/png"
)

type PosterModel struct {
	common   common.Common
	src      string
	image    image.Image
	loaded   bool
	skeleton skeleton.SkeletonModel
}

type PosterMsg = struct {
	src   string
	image image.Image
}

func getSrc(src string) tea.Cmd {

	return func() tea.Msg {
		// client.Logger = &noopLogger{}
		client := retryablehttp.NewClient()
		client.Logger = nil
		resp, err := client.Get(src)

		if err != nil {
			return nil
		}

		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)

		if err != nil {
			return nil
		}

		return PosterMsg{src, img}
	}
}

// The image pixel width is 1/2 of common.Width
// The block characters used to render an image are exactly 3 characters wide
func New(common common.Common, src string) *PosterModel {
	// common.Width = 6
	// common.Height = 9
	errImg := image.NewRGBA(image.Rect(0, 0, 1, 1))
	errImg.Set(0, 0, color.RGBA{252, 52, 2, 0xff})

	skeleton := skeleton.New(common)

	m := &PosterModel{
		src:      src,
		common:   common,
		image:    errImg,
		skeleton: *skeleton,
	}
	// m.SetSize(common.Width, common.Height)
	return m
}

func (m *PosterModel) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height
	// if not loaded ,set skeleton size
	// wm, hm := m.getMargins()

}

func (m *PosterModel) getMargins() (wm, hm int) {
	wm = 0
	hm = 0

	return
}

func (m *PosterModel) Init() tea.Cmd {
	return tea.Batch(getSrc(m.src), m.skeleton.Tick)
}

func (m *PosterModel) Update(msg tea.Msg) (*PosterModel, tea.Cmd) {
	switch msg := msg.(type) {
	case PosterMsg:
		if msg.src == m.src {
			m.image = msg.image
			m.loaded = true
		}
		// TODO check image id
	}

	var cmd tea.Cmd
	_, cmd = m.skeleton.Update(msg)

	return m, cmd
}

func (m *PosterModel) View() string {

	if !m.loaded {
		return m.skeleton.View()
	}

	base := lipgloss.NewStyle() //.Inline(true)
	// base := lipgloss.NewStyle().Inline(true)

	dst := image.NewRGBA(image.Rect(0, 0, m.common.Width, m.common.Height))

	draw.CatmullRom.Scale(dst, dst.Rect, m.image, m.image.Bounds(), draw.Over, nil)

	view := ""

	for y := dst.Bounds().Min.Y; y < dst.Bounds().Max.Y; y++ {

		for x := dst.Bounds().Min.X; x < dst.Bounds().Max.X; x++ {
			r, g, b, _ := dst.At(x, y).RGBA()

			// colors are on a scale from 0 - 65535
			r = r >> 8
			g = g >> 8
			b = b >> 8

			// view += fmt.Sprintf("%v, %v, %v", r, g, b)
			pixel := base.Copy().Background(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b)))

			view += pixel.Render(" ")

		}

		if y == dst.Bounds().Max.Y-1 {
			break
		}

		view += "\n"
	}
	// view += m.skeleton.View()

	return view
}
