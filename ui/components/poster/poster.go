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

type Model struct {
	props    common.Props
	src      string
	image    image.Image
	scaled   *image.RGBA
	loaded   bool
	skeleton *skeleton.Model
}

type PosterMsg = struct {
	src   string
	image image.Image
}

func getSrcCmd(client *retryablehttp.Client, src string) tea.Cmd {

	return func() tea.Msg {
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
func New(p common.Props, src string) *Model {
	errImg := image.NewRGBA(image.Rect(0, 0, 1, 1))
	errImg.Set(0, 0, color.RGBA{252, 52, 2, 0xff})

	m := &Model{
		src:      src,
		props:    p,
		image:    errImg,
		skeleton: skeleton.New(p),
	}
	return m
}

func (m *Model) SetSize(width, height int) {
	m.props.Width = width
	m.props.Height = height

	m.skeleton.SetSize(width, height)
}

func (m *Model) Height() int {
	return m.props.Height
}

func (m *Model) Width() int {
	return m.props.Width
}

type Init struct{}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Init:
		return m, tea.Batch(getSrcCmd(m.props.Global.HttpClient, m.src), m.skeleton.Tick)
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

const (
	top = "▀"
	// bot  = "▄"
	// full = "█"
)

func (m *Model) View() string {

	if !m.loaded {
		return m.skeleton.View()
	}

	view := ""
	if m.scaled == nil ||
		m.scaled.Bounds().Max.X != m.props.Width ||
		m.scaled.Bounds().Max.Y != m.props.Height {

		m.scaled = image.NewRGBA(image.Rect(0, 0, m.props.Width, m.props.Height*2))
		draw.CatmullRom.Scale(m.scaled, m.scaled.Rect, m.image, m.image.Bounds(), draw.Over, nil)
	}

	const text = "Spider-Man: Across the Spider-Verse "
	index := 0
	for y := m.scaled.Bounds().Min.Y; y < m.scaled.Bounds().Max.Y; y += 2 {

		for x := m.scaled.Bounds().Min.X; x < m.scaled.Bounds().Max.X; x++ {
			var fc, bc lipgloss.Color
			{
				r, g, b, _ := m.scaled.At(x, y).RGBA()

				// colors are on a scale from 0 - 65535
				r = r >> 8
				g = g >> 8
				b = b >> 8

				// view += fmt.Sprintf("%v, %v, %v", r, g, b)
				fc = lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b))
			}
			{
				r, g, b, _ := m.scaled.At(x, y+1).RGBA()

				// colors are on a scale from 0 - 65535
				r = r >> 8
				g = g >> 8
				b = b >> 8

				// view += fmt.Sprintf("%v, %v, %v", r, g, b)
				bc = lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b))
			}
			// color := lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b))
			pixel := lipgloss.NewStyle().Foreground(fc).Background(bc)

			// view += pixel.Render(string(text[index%len(text)]))
			view += pixel.Render(top)
			index++

		}

		if y == m.scaled.Bounds().Max.Y-1 {
			break
		}

		view += "\n"
	}
	// view += m.skeleton.View()

	return view
}
