package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/keymap"
)

const (
	host = "0.0.0.0"
	port = 2345
)

func RunServer() {

	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/server_ed25519"),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			lm.Middleware(),
		),
	)

	if err != nil {
		log.Error("server didn't start", "err", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Starting SSH server", "host", host, "port", port)

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("could not start server", "error", err)
		}
	}()

	<-done
	log.Info("Stopping SSH server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal, skipping")
		return nil, nil
	}

	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil

	c := common.Props{
		Global: common.Global{
			AuthState: &common.AuthState{
				Authed: false,
			},
			Config: common.Config{
				TMDB_API_KEY: tmdbKey,
			},

			ReviewMap:  map[int]common.Review{},
			FilmCache:  common.Cache[common.Film]{},
			KeyMap:     keymap.DefaultKeyMap(),
			HttpClient: httpClient,
		},
	}

	p := tea.NewProgram(ui.New(c), tea.WithAltScreen())
	// m := model{
	// 	term:   pty.Term,
	// 	width:  pty.Window.Width,
	// 	height: pty.Window.Height,
	// }
	return nil, []tea.ProgramOption{tea.WithAltScreen()}
}
