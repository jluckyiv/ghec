package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type state int

const (
	loading state = iota
	quitting
	stopped
)

type model struct {
	err     error
	spinner spinner.Model
	state   state
	height  int
	width   int
}

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

var stopKeys = key.NewBinding(
	key.WithKeys(" ", "enter"),
	key.WithHelp("", "press space or enter to stop"),
)

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s, state: loading}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.state = quitting
			return m, tea.Quit
		}
		if key.Matches(msg, stopKeys) {
			if m.state == stopped {
				m.state = loading
				return m, m.spinner.Tick
			}
			if m.state == loading {
				m.state = stopped
			}
			return m, nil
		}
		return m, nil
	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		if m.state == loading {
			m.spinner, cmd = m.spinner.Update(msg)
		}
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	spinnerView := m.spinner.View()
	str := fmt.Sprintf(
		"\n\n   %s Loading forever... %s\n\n    %d, %d",
		spinnerView,
		quitKeys.Help().Desc,
		m.height,
		m.width,
	)
	if m.state == quitting {
		return str + "\n"
	}
	return str
}

func Run() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
