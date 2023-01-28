package image

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

type ImageModel struct {
	common   common.Common
	src      string
	image    image.Image
	loaded   bool
	skeleton skeleton.SkeletonModel
	init     bool
}

type ImageMsg = struct {
	src   string
	image image.Image
}

func getSrc(src string) tea.Cmd {
	// errImg := image.NewRGBA(image.Rect(0, 0, 1, 1))
	// errImg.Set(0, 0, color.RGBA{252, 52, 2, 0xff})

	return func() tea.Msg {
		resp, err := retryablehttp.Get(src)

		if err != nil {
			return nil
		}

		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)

		if err != nil {
			return nil
		}

		return ImageMsg{src, img}
	}
}

func New(common common.Common, src string) *ImageModel {
	common.Width = 10
	common.Height = 15
	errImg := image.NewRGBA(image.Rect(0, 0, 1, 1))
	errImg.Set(0, 0, color.RGBA{252, 52, 2, 0xff})

	skeleton := skeleton.New(common)

	m := &ImageModel{
		src:      src,
		common:   common,
		image:    errImg,
		skeleton: *skeleton,
		init:     false,
	}
	// m.SetSize(common.Width, common.Height)
	return m
}

func (m *ImageModel) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height
	// if not loaded ,set skeleton size
	// wm, hm := m.getMargins()

}

func (m *ImageModel) getMargins() (wm, hm int) {
	wm = 0
	hm = 0

	return
}

func (m *ImageModel) Init() tea.Cmd {
	return tea.Batch(getSrc(m.src), m.skeleton.Tick)
}

func (m *ImageModel) Update(msg tea.Msg) (*ImageModel, tea.Cmd) {
	if !m.init {
		m.init = true
		return m, m.Init()
	}

	switch msg := msg.(type) {
	case ImageMsg:
		if msg.src == m.src {
			m.image = msg.image
			m.loaded = true
		}
		// TODO check image id
	}

	var cmd tea.Cmd
	m.skeleton, cmd = m.skeleton.Update(msg)

	return m, cmd
}

func (m *ImageModel) View() string {

	if !m.loaded {
		return m.skeleton.View()
	}

	base := lipgloss.NewStyle().Inline(true)

	dst := image.NewRGBA(image.Rect(0, 0, m.common.Width, m.common.Height))

	// out, _ := os.Create("dst.png")
	// png.Encode(out, m.image)

	draw.CatmullRom.Scale(dst, dst.Rect, m.image, m.image.Bounds(), draw.Over, nil)

	// out2, _ := os.Create("dst2.png")
	// png.Encode(out2, dst)

	view := ""

	for y := dst.Bounds().Min.Y; y < dst.Bounds().Max.Y; y++ {

		for x := dst.Bounds().Min.X; x < dst.Bounds().Max.X; x++ {
			r, g, b, _ := dst.At(x, y).RGBA()

			r = r >> 8
			g = g >> 8
			b = b >> 8

			// view += fmt.Sprintf("%v, %v, %v", r, g, b)
			pixel := base.Copy().Background(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b)))

			view += pixel.Render("　")

		}

		view += "\n"
	}

	return view
}
